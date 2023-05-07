package mempool

import (
	"financeMgr/src/common/utils"
	"fmt"
	"sync/atomic"
)

type allocObjectFunc func(allocParams interface{}) interface{}

type MemPool struct {
	cap             int
	AllocCount      uint64
	FreeCount       uint64
	objectChan      chan interface{}
	allocObjectFunc allocObjectFunc
	allocParams     interface{}
}

func NewMemPool(cap int, allocObjectFunc allocObjectFunc,
	allocParms interface{}) *MemPool {

	p := &MemPool{}
	p.cap = cap
	p.objectChan = make(chan interface{}, cap)
	p.allocObjectFunc = allocObjectFunc
	p.allocParams = allocParms

	for i := 0; i < cap; i++ {
		obj := p.allocObject()
		p.objectChan <- obj
	}

	return p
}

func (p *MemPool) Alloc() interface{} {
	obj := <-p.objectChan
	utils.Assert(obj != nil)
	atomic.AddUint64(&p.AllocCount, 1)
	return obj
}

func (p *MemPool) AtomicAlloc() interface{} {
	select {
	case obj := <-p.objectChan:
		utils.Assert(obj != nil)
		atomic.AddUint64(&p.AllocCount, 1)
		return obj
	default:
		return nil
	}
}

func (p *MemPool) Free(obj interface{}) {
	utils.Assert(obj != nil)
	select {
	case p.objectChan <- obj:
	default:
		utils.Assert(false)
	}
	atomic.AddUint64(&p.FreeCount, 1)
}

func (p *MemPool) allocObject() interface{} {
	return p.allocObjectFunc(p.allocParams)
}

func (p *MemPool) HasFreeMem() bool {
	return p.AllocCount-p.FreeCount < uint64(p.cap)
}

func FmtPoolCntString(pool *MemPool) (out string) {
	outstring := fmt.Sprintf("\t%20d\t%20d\t%10d ",
		pool.AllocCount, pool.FreeCount,
		pool.AllocCount-pool.FreeCount)
	return outstring
}
