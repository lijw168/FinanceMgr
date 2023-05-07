package mempool

import (
	"financeMgr/src/common/utils"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type ObjPool struct {
	data    []byte
	dataPtr uintptr

	totalCnt  int
	objSize   int
	totalSize int

	stackLock    *sync.Mutex
	stackTop     int
	freeObjStack []unsafe.Pointer

	waitCond *sync.Cond
	waitCnt  int // cnt of goroutines waiting for allocation

	hisAlloc uint64 // protect by stackLock
	hisFree  uint64
	hisWait  uint64

	eraseDirty bool
}

func NewObjPool(objSize, objCnt int, eraseDirty bool) *ObjPool {
	utils.Assert(objSize > 0 && objCnt > 0)

	b := make([]byte, objSize*objCnt)
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	p := &ObjPool{
		data:       b,
		dataPtr:    bh.Data,
		totalCnt:   objCnt,
		objSize:    objSize,
		totalSize:  objCnt * objSize,
		eraseDirty: eraseDirty,
	}
	p.stackLock = &sync.Mutex{}
	p.waitCond = sync.NewCond(p.stackLock)

	// init stack
	p.freeObjStack = make([]unsafe.Pointer, objCnt, objCnt)
	for i := objCnt - 1; i >= 0; i-- {
		ptr := p.dataPtr + uintptr(i*p.objSize)
		objPtr := unsafe.Pointer(ptr)

		p.freeObjStack[p.stackTop] = objPtr
		p.stackTop++
	}

	return p
}

func (p *ObjPool) Alloc() unsafe.Pointer {
	waited := false
	p.stackLock.Lock()
	for p.stackTop == 0 {
		// add waitCnt once
		if !waited {
			waited = true
			p.waitCnt++
			p.hisWait++
		}

		p.waitCond.Wait()
	}

	if waited {
		p.waitCnt--
	}
	utils.Assert(p.waitCnt >= 0)

	// pop object from stack
	p.stackTop--
	objPtr := p.freeObjStack[p.stackTop]

	p.hisAlloc++

	p.stackLock.Unlock()

	return objPtr
}

func (p *ObjPool) Free(objPtr unsafe.Pointer) {
	if p.eraseDirty {
		idx := p.ptrToIdx(objPtr)
		off := idx * p.objSize
		utils.FillZero(p.data[off : off+p.objSize])
	}

	// push object to stack
	p.stackLock.Lock()
	p.freeObjStack[p.stackTop] = objPtr
	p.stackTop++

	utils.Assert(p.waitCnt >= 0)
	if p.waitCnt > 0 {
		p.waitCond.Signal()
	}

	p.hisFree++
	p.stackLock.Unlock()
}

func (p *ObjPool) ptrToIdx(objPtr unsafe.Pointer) int {
	delta := int(uintptr(objPtr) - p.dataPtr)

	utils.Assert(delta >= 0 && delta%p.objSize == 0)
	return delta / p.objSize
}

var objPoolTitle = fmt.Sprintf(
	"\t%8s\t%8s\t%8s\t%8s\t%8s\t%8s\t%8s\n",
	"totalSize(MB)", "totalCnt",
	"hisAlloc", "hisFree", "hisWait",
	"curAvail", "curWait")

func (p *ObjPool) Status(withTitle bool) string {
	var title string
	if withTitle {
		title = objPoolTitle
	}

	return fmt.Sprintf(
		"%s\t%8d\t%8d\t%8d\t%8d\t%8d\t%8d\t%8d",
		title, p.totalSize>>10, p.totalCnt,
		p.hisAlloc, p.hisFree, p.hisWait,
		p.totalCnt-int(p.hisAlloc-p.hisFree), p.waitCnt)
}
