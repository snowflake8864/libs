package list

import (
	//"bufio"
	//	"fmt"
	"github.com/ngaut/log"
	//"os"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Fail()
	}
}

func TestNewListHead(t *testing.T) {
	head := NewListHead()
	log.Debugf("head:%+v", head)
	assert(t, head != nil)
}
func TestAddList(t *testing.T) {

	head := NewListHead()
	head.Init()
	assert(t, head != nil)
	for i := 0; i < 10; i++ {
		n := NewListHead()
		n.Value = i
		head.Add(n)
		head.Value = i
	}
	log.Debugf("head:%+v", head)
	assert(t, head.Value.(int) == 9)
}

func TestAddEntry(t *testing.T) {
	type item struct {
		list  ListHead
		index int
	}
	head := NewListHead()
	var index int
	for i := 0; i < 9; i++ {
		entry := &item{index: index + i}
		entry.list.Value = entry
		head.Add(&entry.list)
	}
	//	head.ForEach(func(list *ListHead) bool {
	head.ForEachSafe(func(list *ListHead) bool {
		log.Debugf("item index is %d", list.Value.(*item).index)
		return true
	})
}
