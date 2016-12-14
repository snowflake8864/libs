package list

import (

//"github.com/ngaut/log"

)

type ListHead struct {
	next, prev *ListHead
	Value      interface{}
}

func NewListHead() *ListHead {
	return new(ListHead).Init()
}

//init list head
func (list *ListHead) Init() *ListHead {
	list.next = list
	list.prev = list
	return list
}

func add(n, prev, next *ListHead) {
	next.prev = n
	n.next = next
	n.prev = prev
	prev.next = n
}

func Add1(head, n *ListHead) {
	add(n, head, head.next)
}

func (head *ListHead) Add(n *ListHead) {
	add(n, head, head.next)
}

func (head *ListHead) AddTail(n *ListHead) {

	add(n, head.prev, head)
}

func del(prev, next *ListHead) {
	next.prev = prev
	prev.next = next
}

func (entry *ListHead) delEntry() {
	del(entry.prev, entry.next)
}

func (entry *ListHead) DelEntry() {
	entry.delEntry()
	entry.next = nil
	entry.prev = nil
}

func (old *ListHead) Replace(n *ListHead) {
	n.next = old.next
	n.next.prev = n
	n.prev = old.prev
	n.prev.next = n
}

func (old *ListHead) ReplaceInit(n *ListHead) {
	old.Replace(n)
	old.Init()
}

func (list *ListHead) DelInit() {
	list.delEntry()
	list.Init()
}

func (head *ListHead) Move(list *ListHead) {
	list.delEntry()
	head.Add(list)
}

func (head *ListHead) MoveTail(list *ListHead) {
	list.delEntry()
	head.AddTail(list)
}

func (head *ListHead) IsLast(list *ListHead) bool {
	return list.next == head
}

func (head *ListHead) Empty() bool {
	return head.next == head
}

func (head *ListHead) EmptyCareful() bool {
	return (head.next == head && head.next == head.prev)
}

func (list *ListHead) Entry() interface{} {
	return list.Value
}

func (head *ListHead) FirstEntry() interface{} {
	return head.next.Entry()
}

func (head *ListHead) LastEntry() interface{} {
	return head.prev.Entry()
}

func (head *ListHead) FirstEntryORNil() interface{} {
	if head.Empty() == true {
		return head.next.Entry()
	} else {
		return nil
	}
}

func (list *ListHead) NextEntry() interface{} {
	return list.next.Entry()
}

func (list *ListHead) PrevEntry() interface{} {
	return list.prev.Entry()
}

type ListFunc func(list *ListHead) bool
type DListFunc func(list1, list2 *ListHead) bool

func (head *ListHead) ForEach(fn ListFunc) {
	for pos := head.next; pos != head; pos = pos.next {
		if !fn(pos) {
			break
		}
	}
}
func (head *ListHead) DForEach(fn DListFunc) {
	for pos0 := head.next; pos0 != head; pos0 = pos0.next {

		for pos1 := head.next; pos1 != head; pos1 = pos1.next {
			if !fn(pos0, pos1) {
				break
			}
		}
	}
}
func (head *ListHead) ForEachPrev(fn ListFunc) {
	for pos := head.prev; ; pos = pos.prev {
		if !fn(pos) {
			break
		}
		if pos == head {
			break
		}
	}
}

func (head *ListHead) ForEachSafe(fn ListFunc) {
	var pos, n *ListHead
	pos = head.next
	for {
		n = pos.next
		if !fn(pos) {
			break
		}

		if pos == head {
			break
		}
		pos = n

	}
}

func (head *ListHead) ForEachPrevSafe(fn ListFunc) {
	var pos, n *ListHead

	pos = head.prev
	for {
		n = pos.prev
		if !fn(pos) {
			break
		}
		if pos == head {
			break
		}
		pos = n

	}

}
