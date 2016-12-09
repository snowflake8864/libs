package list

import (

//"github.com/ngaut/log"

)

type ListHead struct {
	next, prev *ListHead
	Value      interface{}
}

func NewListHead() *ListHead {
	head := &ListHead{}
	head.Init()
	return head
}

//init list head
func (list *ListHead) Init() {
	list.next = list
	list.prev = list
}

func add(n, prev, next *ListHead) {
	next.prev = n
	n.next = next
	n.prev = prev
	prev.next = n
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

func (head *ListHead) ForEach(fn ListFunc) {
	var pos *ListHead
	for pos = head.next; pos != head; pos = pos.next {
		if fn(pos) {
			break
		}
	}
}

func (head *ListHead) ForEachPrev(fn ListFunc) {
	var pos *ListHead
	for pos = head.prev; pos != head; pos = pos.prev {
		if !fn(pos) {
			break
		}
	}
}

func (head *ListHead) ForEachSafe(fn ListFunc) {
	var pos, n *ListHead
	pos = head.next
	for pos != head {
		n = pos.next
		if !fn(pos) {
			break
		}
		pos = n

	}
}

func (head *ListHead) ForEachPrevSafe(fn ListFunc) {
	var pos, n *ListHead
	//for pos = head.prev, n = pos.prev; pos != head; pos = n, n = pos.prev {
	for pos, n = head.prev, pos.prev; pos != head; pos, n = n, pos.prev {
		if !fn(pos) {
			break
		}
	}

	pos = head.prev
	for pos != head {
		n = pos.prev
		if !fn(pos) {
			break
		}
		pos = n

	}

}
