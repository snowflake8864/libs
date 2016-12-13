package acbm

import (
	//   "github.com/snowflake8864/libs/list"
	"container/list"
	"fmt"
	"github.com/ngaut/log"
	"strings"
	//	"bytes"
)

const (
	PATTERN_LEN = 200
)

var CASE_SENSITIVE bool = false

//pattern tree node
type patternTreeNode struct {
	label     int
	depth     int
	ch        byte
	goodShift int
	badShift  int
	oneChild  byte
	childs    [256]*patternTreeNode
	nchild    int
	parent    *patternTreeNode
	pattern   *PatternInfo
}

//
type PatternInfo struct {
	pattern string
	lenth   int
	id      int
	Value   interface{}
}

type Tree struct {
	root           *patternTreeNode
	maxDepth       int
	minPatternSize int
	badShift       [256]int
	//patternList    *PatternInfo
	patternArray []PatternInfo
	patternCount int
	plist        *list.List
	handler      HitHandler
}

type HitHandler interface {
	process(text []byte, data interface{})
	//    mu      sync.RWMutex

}

func New(handler HitHandler) *Tree {
	tree := new(Tree)
	tree.handler = handler
	tree.plist = list.New()
	return tree
	//return &Tree{handler: handler}
}

func toLower(b byte) byte {
	if b >= 65 && b < 90 {
		return b + 32
	} else {
		return b
	}
}

//func (tree *Tree) BuildOne()

func (tree *Tree) Build(dictionary []string) bool {

	arraySize := len(dictionary)
	tree.patternArray = make([]PatternInfo, arraySize)
	for i, s := range dictionary {
		tree.patternArray[i].pattern = s
		tree.patternArray[i].id = i
	}
	var parent *patternTreeNode
	maxPatternLen, minPatternLen := 0, PATTERN_LEN
	root := new(patternTreeNode)
	if root == nil {
		log.Debugf("new root node fail")
		return false
	}
	root.label = -2 //tree root lable
	root.depth = 0  //the depth of tree
	//process string, add these to tree
	npattern := len(tree.patternArray)
	for i := 0; i < npattern; i++ {
		patLen := len(tree.patternArray[i].pattern)
		if patLen == 0 {
			continue
		} else {
			if patLen > PATTERN_LEN {
				patLen = PATTERN_LEN
			}

			if patLen > maxPatternLen {
				maxPatternLen = patLen
			}
			if patLen < minPatternLen {
				minPatternLen = patLen
			}
			parent = root
			var ch_i int
			for ch_i = 0; ch_i < patLen; ch_i++ {
				ch := tree.patternArray[i].pattern[ch_i]
				if !CASE_SENSITIVE {
					ch = toLower(ch)
				}
				if parent.childs[ch] == nil {
					break
				}
				parent = parent.childs[ch]
			}

			if ch_i < patLen {

				for ; ch_i < patLen; ch_i++ {
					ch := tree.patternArray[i].pattern[ch_i]

					if !CASE_SENSITIVE {
						ch = toLower(ch)
					}
					node := new(patternTreeNode)
					if node == nil {
						goto fail
					}
					node.depth = ch_i + 1
					node.ch = ch
					node.label = -1
					//add new node to parent node  with the index
					parent.childs[ch] = node
					if !CASE_SENSITIVE {
						if (ch >= 'a') && (ch <= 'z') {
							parent.childs[ch-32] = node
						}
					}
					parent.nchild++
					parent.oneChild = ch
					node.parent = parent
					parent = node
				}
			}
		}
		parent.label = i
		parent.pattern = &tree.patternArray[i]
		//log.Debugf("pattern lable %d", i)
		//		tree.plist.PushFront(pData)
	}
	tree.patternCount = npattern
	tree.root = root
	tree.maxDepth = maxPatternLen
	tree.minPatternSize = minPatternLen
	log.Debugf("Build pattern ok")
	return true
fail:
	if tree.root != nil {
		cleanTree(tree.root)
		tree.root = nil
	}
	return false
}

func cleanTree(root *patternTreeNode) {
	for i := 0; i < 256; i++ {
		if root.childs[i] != nil {
			cleanTree(root.childs[i])
			root.childs[i] = nil
		}
	}
	return
}

func (tree *Tree) Clean() {
	if tree.root != nil {
		cleanTree(tree.root)
		tree.root = nil
	}
	return
}

func (tree *Tree) ComputeBCShifts() {
	for i := 0; i < 256; i++ {
		tree.badShift[i] = tree.minPatternSize
	}
	for i := tree.minPatternSize - 1; i > 0; i-- {
		for j := 0; j < tree.patternCount; j++ {
			ch := tree.patternArray[j].pattern[i]
			if !CASE_SENSITIVE {
				ch = toLower(ch)
			}

			//fmt.Println("ch", ch)
			tree.badShift[ch] = i
			if !CASE_SENSITIVE {
				if (ch > 'a') && (ch < 'z') {
					tree.badShift[ch-32] = i
				}
			}
		}
	}
	return
}

func _initGSshifts(root *patternTreeNode, shift int) {
	if root.label != -2 {
		root.goodShift = shift
	}

	for i := 0; i < 256; i++ {

		if !CASE_SENSITIVE {
			if (i > 'A') && (i < 'Z') {
				continue
			}
		}
		if root.childs[i] != nil {
			_initGSshifts(root.childs[i], shift)
		}
	}
}
func (tree *Tree) initGSshifts() {
	_initGSshifts(tree.root, tree.minPatternSize)
}

func (tree *Tree) setGSshift(pat string, depth, shift int) int {
	if tree == nil || tree.root == nil {
		return -1
	}
	node := tree.root
	for i := 0; i < depth; i++ {

		node = node.childs[pat[i]]
		if node == nil {

			return -1
		}
	}

	// get the little offset

	if node.goodShift >= shift {
		node.goodShift = shift
	}
	//fmt.Println("shift:", node.goodShift)
	return 0
}

func (tree *Tree) computeGSShift(pat1 string, pat2 string) bool {

	pat1Len := len(pat1)
	pat2Len := len(pat2)
	if pat1Len < 0 || pat2Len < 0 {
		return false
	}
	if pat1Len == 0 || pat2Len == 0 {
		return true
	}

	if !CASE_SENSITIVE {
		pat1 = strings.ToLower(pat1)
		pat2 = strings.ToLower(pat2)
	}
	firstChar := pat1[0]
	i := 0
	for i = 1; i < pat2Len; i++ {

		if pat2[i] != firstChar {
			break
		}
	}

	tree.setGSshift(pat1, 1, i)

	i = 1
	for {
		for i < pat2Len && pat2[i] != firstChar {
			i++
		}
		if i == pat2Len {
			break
		}
		pat2Index := i
		pat1Index := 0
		offset := i

		if offset > tree.minPatternSize {
			break
		}

		for pat2Index < pat2Len && pat1Index < pat1Len {

			if pat1[pat1Index] != pat2[pat2Index] {

				break
			}
			pat1Index++ //是比较位的字符的深度
			pat2Index++
		}
		if pat2Index == pat2Len { // 关键字pat1前缀是关键字pat2后缀

			for j := pat1Index; j < pat1Len; j++ {
				tree.setGSshift(pat1, j+1, offset)
			}
		} else { // pat1的前缀是pat2的子串,被pat2包含

			tree.setGSshift(pat1, pat1Index+1, offset) //字符所在的深度和序号差一位
		}
		i++
	}
	return true
}

func (tree *Tree) computeGSShifts() {
	for pat_i := 0; pat_i < tree.patternCount; pat_i++ {
		for pat_j := 0; pat_j < tree.patternCount; pat_j++ {
			ppat_i := tree.patternArray[pat_i].pattern
			ppat_j := tree.patternArray[pat_j].pattern
			tree.computeGSShift(ppat_i, ppat_j)
		}
	}
}

func (tree *Tree) ComputeShifts() {
	tree.ComputeBCShifts()
	tree.initGSshifts()
	tree.computeGSShifts()
}

func (tree *Tree) Search(text []byte) int {
	//	log.Debugf("text:%s", text)
	nmatched := 0

	textLen := len(text)
	if textLen < tree.minPatternSize {
		return nmatched
	}

	baseIndex := textLen - tree.minPatternSize
	for baseIndex >= 0 {

		curIndex := baseIndex
		node := tree.root

		//log.Debugf("[%c]", text[curIndex])
		for node.childs[text[curIndex]] != nil {

			log.Debugf("%c", text[curIndex])
			node = node.childs[text[curIndex]]

			if node.label >= 0 {

				log.Debugf("Matched(%d) ", node.label)
				log.Debugf("%s ", text[baseIndex:])

				tree.handler.process(text[baseIndex:], node.pattern)
				nmatched++
			}
			curIndex++
			if curIndex >= textLen {

				break
			}
		}

		if node.nchild > 0 {

			gsShift := node.childs[node.oneChild].goodShift
			bcShift := 0
			if curIndex < textLen {
				bcShift = tree.badShift[text[curIndex]] + baseIndex - curIndex

			} else {
				bcShift = 1
			}
			realShift := bcShift
			if gsShift > bcShift {
				realShift = gsShift
			}
			baseIndex -= realShift
		} else {
			baseIndex--
		}
	}

	return nmatched
}

func (root *patternTreeNode) print() {
	if root == nil {
		return
	}
	fmt.Print("ch:", root.ch, " goodShift:", root.goodShift, " lable:", root.label, "depth:", root.depth, " childs:")

	for i := 0; i < 256; i++ {
		if !CASE_SENSITIVE {
			if (i >= 'A') && (i <= 'Z') {
				continue
			}
		}
		if root.childs[i] != nil {
			fmt.Print(byte(root.childs[i].ch), " ")
		}
	}
	fmt.Println()

	for i := 0; i < 256; i++ {
		if root.childs[i] != nil {
			root.childs[i].print()
		}
	}
	return
}

func (tree *Tree) Print() {
	fmt.Println("-----ACTree information----")
	if tree.root != nil {
		tree.root.print()

	}
	return
}
