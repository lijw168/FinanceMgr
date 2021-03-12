package mempool

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"
)

type Ctx struct {
	groupID  string
	ref      int
	lock     sync.Mutex
	entryMap map[uint64]string
	buf      [512]byte
}

func TestObjPool_Basic(t *testing.T) {
	pool := NewObjPool(int(unsafe.Sizeof(Ctx{})), 10, true)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 50; i++ {
		allocCnt := r.Intn(11) // allocCnt is in [0,10]
		useObj(t, pool, allocCnt)
	}
}

func TestObjPool_ObjAddress(t *testing.T) {
	// sanity check of object address, increased by object size
	ctxSize := unsafe.Sizeof(Ctx{})
	pool := NewObjPool(int(ctxSize), 10, true)
	objPtrs := make([]unsafe.Pointer, 0)

	firstObjPtr := pool.Alloc()
	preObjAddr := uintptr(firstObjPtr)
	objPtrs = append(objPtrs, firstObjPtr)

	for i := 1; i < 10; i++ {
		curObjPtr := pool.Alloc()
		curObjAddr := uintptr(curObjPtr)
		objPtrs = append(objPtrs, curObjPtr)

		assert.Equal(t, curObjAddr, preObjAddr+ctxSize)

		preObjAddr = curObjAddr
	}

	for _, objPtr := range objPtrs {
		pool.Free(objPtr)
	}
}

func TestObjPool_FreeList(t *testing.T) {
	pool := NewObjPool(int(unsafe.Sizeof(Ctx{})), 10, true)
	objPtrs := make([]unsafe.Pointer, 0)

	// exhaust object pool
	for i := 0; i < 10; i++ {
		objPtr := pool.Alloc()
		objPtrs = append(objPtrs, objPtr)
	}

	// expect stack is emtpy
	assert.True(t, pool.stackTop == 0)

	// free obj-7, obj-8, obj-9, expect stack: [9, 8, 7]
	pool.Free(objPtrs[7])
	pool.Free(objPtrs[8])
	pool.Free(objPtrs[9])
	assert.True(t, pool.stackTop == 3)
	assert.True(t, pool.ptrToIdx(pool.freeObjStack[2]) == 9)
	assert.True(t, pool.ptrToIdx(pool.freeObjStack[1]) == 8)
	assert.True(t, pool.ptrToIdx(pool.freeObjStack[0]) == 7)
	objPtrs = objPtrs[:7]

	// free the left
	for _, objPtr := range objPtrs {
		pool.Free(objPtr)
	}
}

func TestObjPool_WaitForRelease(t *testing.T) {
	var wg sync.WaitGroup
	pool := NewObjPool(int(unsafe.Sizeof(Ctx{})), 10, true)

	useSingleObjFunc := func() {
		useSingleObj(t, pool)
		wg.Done()
	}

	wg.Add(50)
	for i := 0; i < 50; i++ {
		go useSingleObjFunc()
	}

	wg.Wait()
}

func useSingleObj(t *testing.T, pool *ObjPool) {
	useObj(t, pool, 1)
}

func useObj(t *testing.T, pool *ObjPool, allocCnt int) {
	objPtrs := make([]unsafe.Pointer, 0)

	for i := 0; i < allocCnt; i++ {
		objPtr := pool.Alloc()
		objPtrs = append(objPtrs, objPtr)

		ctx := (*Ctx)(objPtr)

		// expect object empty
		assert.Equal(t, ctx.groupID, "")
		assert.Equal(t, ctx.ref, 0)
		assert.True(t, len(ctx.entryMap) == 0)
		assert.True(t, ctx.buf[0] == 0)

		// use object
		ctx.groupID = "id"
		ctx.ref = 2
		ctx.lock.Lock()
		ctx.lock.Unlock()
		ctx.entryMap = make(map[uint64]string)
		ctx.entryMap[0] = "map"
		ctx.buf[0] = 1

		// sleep some time to simulate using object
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		ms := r.Intn(20)
		time.Sleep(time.Millisecond * time.Duration(ms))
	}

	for _, objPtr := range objPtrs {
		pool.Free(objPtr)
	}
}
