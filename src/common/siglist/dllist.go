package siglist

import (
	"common/utils"
	"sync"
)

type List_head struct {
	Next *List_head
	Prev *List_head
}

type Dllist struct {
	//Must In First place
	List List_head

	WaitCh chan struct{}
	len    int
	Lock   sync.Mutex
}

func (dl *Dllist) list_add(entry *List_head, prev *List_head,
	next *List_head) {
	next.Prev = entry
	entry.Next = next
	entry.Prev = prev
	prev.Next = entry
}

func (dl *Dllist) list_add_tail(entry *List_head, head *List_head) {
	dl.list_add(entry, head.Prev, head)
}

func (dl *Dllist) list_add_head(entry *List_head, head *List_head) {
	dl.list_add(entry, head, head.Next)
}

func (dl *Dllist) list_move_to_head(entry *List_head, head *List_head) {
	dl.list_del(entry)
	dl.list_add_head(entry, head)
}

func (dl *Dllist) Init() {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()

	dl.WaitCh = make(chan struct{}, 1)
	dl.List.Prev = &dl.List
	dl.List.Next = &dl.List
	dl.len = 0
}

func (dl *Dllist) PushBack(entry *List_head, head *List_head) {
	dl.Lock.Lock()
	dl.list_add_tail(entry, head)
	dl.len++
	dl.Lock.Unlock()

	dl.WakeUp()
}

func (dl *Dllist) PushFront(entry *List_head, head *List_head) {
	dl.Lock.Lock()
	dl.list_add_head(entry, head)
	dl.len++
	dl.Lock.Unlock()

	dl.WakeUp()
}

func (dl *Dllist) PushBackNoWake(entry *List_head) {
	dl.Lock.Lock()
	dl.list_add_tail(entry, &dl.List)
	dl.len++
	dl.Lock.Unlock()
}

func (dl *Dllist) MoveToFront(entry *List_head) {
	dl.Lock.Lock()
	dl.list_move_to_head(entry, &dl.List)
	dl.Lock.Unlock()

	dl.WakeUp()
}

func (dl *Dllist) WakeUp() {
	select {
	case dl.WaitCh <- struct{}{}:
	default:
	}
}

func (dl *Dllist) _list_del(prev *List_head, next *List_head) {
	next.Prev = prev
	prev.Next = next
}

func (dl *Dllist) list_del(entry *List_head) {
	dl._list_del(entry.Prev, entry.Next)
	entry.Prev = nil
	entry.Next = nil
}

func (dl *Dllist) list_front() *List_head {
	return dl.List.Next
}

func (dl *Dllist) list_back() *List_head {
	return dl.List.Prev
}

func (dl *Dllist) list_empty(head *List_head) bool {
	return head.Next == head
}

func (dl *Dllist) PopFront() *List_head {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()
	if dl.list_empty(&dl.List) == true {
		utils.Assert(dl.len == 0)
		return nil
	}

	utils.Assert(dl.len > 0)
	front := dl.list_front()
	utils.Assert(front != nil)
	utils.Assert(front != &dl.List)
	dl.list_del(front)
	dl.len--
	utils.Assert(dl.len >= 0)
	return front
}

func (dl *Dllist) PopBack() *List_head {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()
	if dl.list_empty(&dl.List) == true {
		utils.Assert(dl.len == 0)
		return nil
	}

	utils.Assert(dl.len > 0)
	back := dl.list_back()
	utils.Assert(back != nil)
	utils.Assert(back != &dl.List)
	dl.list_del(back)
	dl.len--
	utils.Assert(dl.len >= 0)
	return back
}

func (dl *Dllist) GetBack() *List_head {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()
	if dl.list_empty(&dl.List) == true {
		utils.Assert(dl.len == 0)
		return nil
	}

	utils.Assert(dl.len > 0)
	back := dl.list_back()
	utils.Assert(back != nil)
	utils.Assert(back != &dl.List)
	return back
}

func (dl *Dllist) DelEntry(entry *List_head) {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()
	if dl.list_empty(&dl.List) == true {
		utils.Assert(dl.len == 0)
	}

	utils.Assert(dl.len > 0)
	dl.list_del(entry)
	dl.len--
	utils.Assert(dl.len >= 0)
}

func (dl *Dllist) Len() int {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()

	return dl.len
}

func (dl *Dllist) IsHead(entry *List_head) bool {
	dl.Lock.Lock()
	defer dl.Lock.Unlock()

	utils.Assert(entry != nil)
	return (entry == dl.list_front())
}

/*
func TEST_Dllist() {
	var DL_list siglist.Dllist
	DL_list.Init()
	fmt.Printf("%p %p %d\n", DL_list.List.Prev,
		DL_list.List.Next, DL_list.Len())

	fmt.Printf(" ------------------\n")
	entry1 := new(entry)
	fmt.Printf("%p \n", &entry1.node)
	entry1.value = 1
	DL_list.PushBack(&entry1.node, &DL_list.List)

	fmt.Printf("%p %p  %d\n", DL_list.List.Prev,
		DL_list.List.Next, DL_list.Len())
	fmt.Printf(" ------------------\n")

	entry2 := new(entry)
	fmt.Printf("%p \n", &entry2.node)
	entry2.value = 2
	DL_list.PushBack(&entry2.node, &DL_list.List)
	fmt.Printf("%p %p %d\n", DL_list.List.Prev,
		DL_list.List.Next, DL_list.Len())
	fmt.Printf(" ------------------\n")

	entry3 := new(entry)
	fmt.Printf("%p \n", &entry3.node)
	DL_list.PushBack(&entry3.node, &DL_list.List)
	entry3.value = 3
	fmt.Printf("%p %p %d\n", DL_list.List.Prev,
		DL_list.List.Next, DL_list.Len())
	fmt.Printf(" ------------------\n")

	fmt.Printf(" <<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>\n")

	en := DL_list.PopFront()
	en1 := (*entry)(unsafe.Pointer(en))
	fmt.Printf(" en1 %p %d %d\n", en, en1.value, DL_list.Len())

	en = DL_list.PopFront()
	en2 := (*entry)(unsafe.Pointer(en))
	fmt.Printf(" en2 %p %d %d\n", en, en2.value, DL_list.Len())

	en = DL_list.PopFront()
	en3 := (*entry)(unsafe.Pointer(en))
	fmt.Printf(" en3 %p %d %d\n", en, en3.value, DL_list.Len())

	en = DL_list.PopFront()
	fmt.Printf(" en4 %p %d \n", en, DL_list.Len())
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
	fmt.Printf("%p %p \n", DL_list.List.Prev, DL_list.List.Next)
}
*/
