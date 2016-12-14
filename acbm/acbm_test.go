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

/*
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
*/

type pattern struct {
	offset int
	index  int
}

type Mystruct struct{}

func (mystruct *Mystruct) process(text []byte, data interface{}, offset int) {
	//fmt.Printf("%s\n", string(text))
	fmt.Println(string(text[0:]))
	p := data.(*pattern)
	fmt.Printf("Pattern:%d--%d\n", p.index, offset)
}

func TestBuildPattern(t *testing.T) {
	tree := New(&Mystruct{})
	assert(t, tree != nil)
	keys := []string{"snow", "zebra", "AC-BM", "雪花"}
	for i := 0; i < len(keys); i++ {
		p := new(pattern)
		p.index = i

		tree.CreatePattern(keys[i], p)
	}
	tree.BuildPattern()
	tree.ComputeShifts()

	tree.ComputeBCShifts()
	nmatched := tree.Search([]byte("A Golang implementation of the AC-BM string  multiPattern matching algorithm,by zebra, snowflake, Countless stars in the sky, but the moon is only one, 雪花啤酒好喝吗?"))
	assert(t, nmatched == 3)
	//	tree.Print()
}
