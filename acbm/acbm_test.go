package acbm

import (
	//   "regexp"
	//    "strings"
	"fmt"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Fail()
	}
}

type Mystruct struct{}

func (mystruct *Mystruct) process(text []byte, data interface{}) {
	fmt.Println(text)
	Pattern := data.(*PatternInfo)
	fmt.Println("Pattern:%v", Pattern)
}

func TestNewACBM(t *testing.T) {
	tree := New(&Mystruct{})
	assert(t, tree != nil)
}

func TestBuild(t *testing.T) {
	tree := New(&Mystruct{})
	assert(t, tree != nil)
	b := tree.Build([]string{"snow", "zebra", "AC-BM", "雪花"})
	assert(t, b == true)
	//tree.Print()
	tree.ComputeShifts()

	tree.ComputeBCShifts()
	//	tree.Print()
	//	var matchedItems []matcheInfo
	//	matchedItems := make([]matcheInfo, 20)
	//nmatched := tree.Search([]byte("A Golang implementation of the AC-BM string  multiPattern matching algorithm,by zebra, snowflake"), matchedItems, 20)
	nmatched := tree.Search([]byte("A Golang implementation of the AC-BM string  multiPattern matching algorithm,by zebra, snowflake"))
	assert(t, nmatched == 3)
	//tree.Print()
}
