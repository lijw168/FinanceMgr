package siglist

import (
	"container/list"
	"sync"
)

type SignalList struct {
	WaitCh chan struct{}
	List   *list.List
	Lock   sync.RWMutex
}

func (sl *SignalList) Init() {
	sl.WaitCh = make(chan struct{}, 1)
	sl.List = list.New()
}

func (sl *SignalList) PushBack(entry interface{}) {
	sl.Lock.Lock()
	sl.List.PushBack(entry)
	sl.Lock.Unlock()

	sl.WakeUp()
}

func (sl *SignalList) PushFront(entry interface{}) {
	sl.Lock.Lock()
	sl.List.PushFront(entry)
	sl.Lock.Unlock()

	sl.WakeUp()
}

func (sl *SignalList) WakeUp() {
	select {
	case sl.WaitCh <- struct{}{}:
	default:
	}
}

func (sl *SignalList) PopFront() *list.Element {
	sl.Lock.Lock()
	front := sl.List.Front()
	if front != nil {
		sl.List.Remove(front)
	}
	sl.Lock.Unlock()
	return front
}

func (sl *SignalList) Clear() {
	for sl.PopFront() != nil {
	}
}

func (sl *SignalList) IsEmpty() bool {
	return sl.Len() == 0
}

func (sl *SignalList) Len() int {
	sl.Lock.Lock()
	defer sl.Lock.Unlock()
	return sl.List.Len()
}
