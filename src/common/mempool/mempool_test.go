package mempool

import (
	"github.com/stretchr/testify/assert"
	"common/utils"
	"testing"
	"unsafe"
)

type TempObj struct{}

func allocTempObj(allocParams []int) interface{} {
	utils.Assert(len(allocParams) == 0)
	return &TempObj{}
}

func TestAllocObject(t *testing.T) {
	cap := 10

	p, err := NewMemPool(cap, allocTempObj)
	assert.Nil(t, err)

	obj := p.Alloc()
	_, ok := obj.(*TempObj)
	assert.True(t, ok)
	p.Free(obj)
}

func TestAllocAndFree(t *testing.T) {
	cap := 10

	p, err := NewMemPool(cap, allocTempObj)
	assert.Nil(t, err)
	assert.Equal(t, cap, len(p.objectChan))

	// Exhaust capacity
	objs := make([]interface{}, 0, cap)
	for i := 0; i < cap; i++ {
		obj := p.Alloc()
		objs = append(objs, obj)
	}
	assert.Equal(t, 0, len(p.objectChan))
	assert.Equal(t, cap, int(p.allocCount))

	// Free a obj and re-alloc from pool. Expect re-alloc that obj.
	objA := objs[0]
	ptrA := getPtrOfTempObj(objA)
	p.Free(objA)
	objs = objs[1:]

	objB := p.Alloc()
	ptrB := getPtrOfTempObj(objB)
	assert.Equal(t, ptrA, ptrB)
	p.Free(objB)

	// Free allocated object
	for _, obj := range objs {
		p.Free(obj)
	}
	assert.Equal(t, cap, len(p.objectChan))
	assert.Equal(t, 0, int(p.allocCount))
}

func getPtrOfTempObj(obj interface{}) uintptr {
	tempObj, _ := obj.(*TempObj)
	return uintptr(unsafe.Pointer(tempObj))
}
