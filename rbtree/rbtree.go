package rbtree

import (
	"github.com/funny/slab"
	"reflect"
	"unsafe"
)

type Item interface{}

const BLACK = 0
const RED = 1

var rbPool *slab.SyncPool

type Node struct {
	KeyValue   Item /* generic key */
	Parent     *Node
	LeftChild  *Node
	RightChild *Node
	color      uint8
	addr       []byte
}

type CompareFunc func(keyA, keyB Item) int

type RBtree struct {
	Root *Node /* root of the tree */
	//	minNode *Node
	//	maxNode *Node
	//	size    uint
	/* this function directly provides the keys as parameters not the nodes. */
	Compare CompareFunc
}

func Byte2Pointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	)
}

const NODESIZE = int(unsafe.Sizeof(Node{}))

/* generic leaf node */
var LEAF = Node{
	KeyValue:   nil, /* key */
	Parent:     nil,
	LeftChild:  nil,
	RightChild: nil,   /* Parent, left,right */
	color:      BLACK, /* color */
}

func (tree *RBtree) For_each_safe(fn func(node *Node) bool) Item {
	pos := TreeFirstNode(tree.Root)
	for pos != &LEAF {
		if fn(pos) == true {
			return pos.KeyValue
			//break
		}
		pos = TreeNext(pos)
	}
	return nil
}
func (tree *RBtree) OnlyInitTree(compareFunction CompareFunc) {
	tree.Root = &LEAF
	tree.Compare = compareFunction
}

func InitTree(compareFunction CompareFunc) *RBtree {

	tree := &RBtree{
		Root:    &LEAF,
		Compare: compareFunction,
	}
	return tree
}
func IS_LEAF(x *Node) bool {
	return x == &LEAF
}

func IS_NOT_LEAF(x *Node) bool {
	return (x != &LEAF)
}

//const INIT_LEAF = (&LEAF)
var INIT_LEAF *Node = (&LEAF)

func AM_I_LEFT_CHILD(x *Node) bool {
	return (x == x.Parent.LeftChild)
}

func AM_I_RIGHT_CHILD(x *Node) bool {
	return (x == x.Parent.RightChild)
}

func PARENT(x *Node) *Node {
	return (x.Parent)
}

func GRANPA(x *Node) *Node {
	return (x.Parent.Parent)
}

func TreeFirstNode(node *Node) *Node {
	for IS_NOT_LEAF(node.LeftChild) {
		node = node.LeftChild
	}

	return node
}

func (tree *RBtree) FirstKey() Item {
	node := tree.Root
	if IS_NOT_LEAF(node.LeftChild) {
		node = node.LeftChild
	}
	if node == nil {
		return nil
	}
	return node.KeyValue
}

func TreeLastNode(node *Node) *Node {
	for IS_NOT_LEAF(node.RightChild) {
		node = node.RightChild
	}
	return node
}

func TreeNext(node *Node) *Node {

	if node == nil {
		return nil
	}
	if IS_NOT_LEAF(node.RightChild) {
		node = node.RightChild
		for IS_NOT_LEAF(node.LeftChild) {
			node = node.LeftChild
		}
	} else {
		if IS_NOT_LEAF(node.Parent) && AM_I_LEFT_CHILD(node) {
			node = node.Parent
		} else {
			for IS_NOT_LEAF(node.Parent) && AM_I_RIGHT_CHILD(node) {
				node = node.Parent
			}
			node = node.Parent
		}
	}

	return node
}

func TreePrev(node *Node) *Node {

	if IS_NOT_LEAF(node.LeftChild) {
		node = node.LeftChild
		for IS_NOT_LEAF(node.RightChild) {
			node = node.RightChild
		}
	} else {
		if IS_NOT_LEAF(node.Parent) && AM_I_RIGHT_CHILD(node) {
			node = node.Parent
		} else {
			for IS_NOT_LEAF(node) && AM_I_LEFT_CHILD(node) {
				node = node.Parent
			}
			node = node.Parent
		}
	}

	return node
}

func (tree *RBtree) FindNode(key Item) *Node {

	found := tree.Root

	for IS_NOT_LEAF(found) {

		result := tree.Compare(found.KeyValue, key)
		if result == 0 {
			return found
		} else if result < 0 {
			found = found.RightChild
		} else {
			found = found.LeftChild
		}
	}
	return nil
}

func (tree *RBtree) Empty() bool {
	return tree.Root == nil || IS_LEAF(tree.Root)
}

func (tree *RBtree) FindKey(key Item) Item {

	var found *Node

	found = tree.Root
	for IS_NOT_LEAF(found) {
		/*
			if IS_LEAF(found) {
				break
			}
		*/
		result := tree.Compare(found.KeyValue, key)
		if result == 0 {
			return found.KeyValue
		} else if result < 0 {
			found = found.RightChild
		} else {
			found = found.LeftChild
		}

	}
	return nil
}

func (tree *RBtree) Insert(key Item) Item {

	last_node := INIT_LEAF
	temp := tree.Root
	var insert *Node
	var LocalKey Item
	result := int(0)

	if IS_NOT_LEAF(tree.Root) {
		LocalKey = tree.FindKey(key)
	} else {
		LocalKey = nil
	}

	/* if node already in, bail out */
	if LocalKey != nil {
		return LocalKey
	} else {
		insert = tree.createNode(key)

		if insert == nil {
			/* to let the user know that it couldn't insert */
			return Item(&LEAF)
		}
	}

	/* search for the place to insert the new node */
	for IS_NOT_LEAF(temp) {
		last_node = temp
		result = tree.Compare(insert.KeyValue, temp.KeyValue)
		if result < 0 {
			temp = temp.LeftChild
		} else {
			temp = temp.RightChild
		}
	}

	/* make the needed connections */
	insert.Parent = last_node

	if IS_LEAF(last_node) {
		tree.Root = insert
	} else {
		result = tree.Compare(insert.KeyValue, last_node.KeyValue)
		if result < 0 {
			last_node.LeftChild = insert
		} else {
			last_node.RightChild = insert
		}
	}

	/* fix colour issues */
	tree.fixInsertCollisions(insert)
	return nil
}

func (tree *RBtree) Delete(key Item) Item {

	if key == nil {
		return nil
	}
	delete := tree.FindNode(key)

	/* this key isn't in the tree, bail out */
	if delete == nil {
		return nil
	}
	var temp *Node
	lkey := delete.KeyValue
	nodeColor := tree.deleteCheckSwitch(delete, &temp)

	/* deleted node is black, this will mess up the black path property */
	if nodeColor == BLACK {
		tree.fixDeleteCollisions(temp)
	}

	RBNodeFree(delete)
	return lkey
}

func (tree *RBtree) deleteCheckSwitch(delete *Node, temp **Node) uint8 {
	ltemp := delete
	nodeColor := delete.color
	if IS_LEAF(delete.LeftChild) {
		*temp = ltemp.RightChild
		tree.switchNodes(ltemp, ltemp.RightChild)
	} else {
		if IS_LEAF(delete.RightChild) {
			_ltemp := delete
			*temp = _ltemp.LeftChild
			tree.switchNodes(_ltemp, _ltemp.LeftChild)
		} else {
			nodeColor = tree.deleteNode(delete, temp)
		}
	}
	return nodeColor
}

func (tree *RBtree) deleteNode(d *Node, temp **Node) uint8 {

	ltemp := d
	min := TreeFirstNode(d.RightChild)
	nodeColor := min.color

	*temp = min.RightChild
	if min.Parent == ltemp && IS_NOT_LEAF(*temp) {
		(*temp).Parent = min
	} else {
		tree.switchNodes(min, min.RightChild)
		min.RightChild = ltemp.RightChild
		if IS_NOT_LEAF(min.RightChild) {
			min.RightChild.Parent = min
		}
	}

	tree.switchNodes(ltemp, min)
	min.LeftChild = ltemp.LeftChild

	if IS_NOT_LEAF(min.LeftChild) {
		min.LeftChild.Parent = min
	}
	min.color = ltemp.color
	return nodeColor
}

func (tree *RBtree) switchNodes(nodeA *Node, nodeB *Node) {

	if IS_LEAF(nodeA.Parent) {
		tree.Root = nodeB
	} else {
		if IS_NOT_LEAF(nodeA) {
			if AM_I_LEFT_CHILD(nodeA) {
				nodeA.Parent.LeftChild = nodeB
			} else {
				nodeA.Parent.RightChild = nodeB
			}
		}
	}
	if IS_NOT_LEAF(nodeB) {
		nodeB.Parent = nodeA.Parent
	}

}

/*
 * This function fixes the possible collisions in the tree.
 * Eg. if a node is red his children must be black !
 */
func (tree *RBtree) fixInsertCollisions(node *Node) {
	var temp *Node

	for node.Parent.color == RED && IS_NOT_LEAF(GRANPA(node)) {
		if AM_I_RIGHT_CHILD(node.Parent) {
			temp := GRANPA(node).LeftChild
			if temp.color == RED {
				node.Parent.color = BLACK
				temp.color = BLACK
				GRANPA(node).color = RED
				node = GRANPA(node)
			} else if temp.color == BLACK {
				if node == node.Parent.LeftChild {
					node = node.Parent
					rotateToRight(tree, node)
				}

				node.Parent.color = BLACK
				GRANPA(node).color = RED
				rotateToLeft(tree, GRANPA(node))
			}
		} else if AM_I_LEFT_CHILD(node.Parent) {
			temp = GRANPA(node).RightChild
			if temp.color == RED {
				node.Parent.color = BLACK
				temp.color = BLACK
				GRANPA(node).color = RED
				node = GRANPA(node)
			} else if temp.color == BLACK {
				if AM_I_RIGHT_CHILD(node) {
					node = node.Parent
					rotateToLeft(tree, node)
				}

				node.Parent.color = BLACK
				GRANPA(node).color = RED
				rotateToRight(tree, GRANPA(node))
			}
		}
	}
	/* make sure that the Root of the tree stays black */
	tree.Root.color = BLACK
}

func (tree *RBtree) fixDeleteCollisions(node *Node) {
	var temp *Node

	for node != tree.Root && node.color == BLACK && IS_NOT_LEAF(node) {
		if AM_I_LEFT_CHILD(node) {

			temp = node.Parent.RightChild
			if temp.color == RED {
				temp.color = BLACK
				node.Parent.color = RED
				rotateToLeft(tree, node.Parent)
				temp = node.Parent.RightChild
			}

			if temp.LeftChild.color == BLACK && temp.RightChild.color == BLACK {
				temp.color = RED
				node = node.Parent
			} else {
				if temp.RightChild.color == BLACK {
					temp.LeftChild.color = BLACK
					temp.color = RED
					rotateToRight(tree, temp)
					temp = temp.Parent.RightChild
				}

				temp.color = node.Parent.color
				node.Parent.color = BLACK
				temp.RightChild.color = BLACK
				rotateToLeft(tree, node.Parent)
				node = tree.Root
			}
		} else {
			temp = node.Parent.LeftChild
			if temp.color == RED {
				temp.color = BLACK
				node.Parent.color = RED
				rotateToRight(tree, node.Parent)
				temp = node.Parent.LeftChild
			}

			if temp.RightChild.color == BLACK && temp.LeftChild.color == BLACK {
				temp.color = RED
				node = node.Parent
			} else {
				if temp.LeftChild.color == BLACK {
					temp.RightChild.color = BLACK
					temp.color = RED
					rotateToLeft(tree, temp)
					temp = temp.Parent.LeftChild
				}

				temp.color = node.Parent.color
				node.Parent.color = BLACK
				temp.LeftChild.color = BLACK
				rotateToRight(tree, node.Parent)
				node = tree.Root
			}
		}
	}
	node.color = BLACK
}

func rotateToLeft(tree *RBtree, node *Node) {

	temp := node.RightChild

	if temp == &LEAF {
		return
	}
	node.RightChild = temp.LeftChild

	if IS_NOT_LEAF(temp.LeftChild) {
		temp.LeftChild.Parent = node
	}
	temp.Parent = node.Parent

	if IS_LEAF(node.Parent) {
		tree.Root = temp
	} else {
		if node == node.Parent.LeftChild {
			node.Parent.LeftChild = temp
		} else {
			node.Parent.RightChild = temp
		}
	}
	temp.LeftChild = node
	node.Parent = temp
}

func rotateToRight(tree *RBtree, node *Node) {

	temp := node.LeftChild
	node.LeftChild = temp.RightChild

	if temp == &LEAF {
		return
	}

	if IS_NOT_LEAF(temp.RightChild) {
		temp.RightChild.Parent = node
	}

	temp.Parent = node.Parent

	if IS_LEAF(node.Parent) {
		tree.Root = temp
	} else {
		if node == node.Parent.RightChild {
			node.Parent.RightChild = temp
		} else {
			node.Parent.LeftChild = temp
		}
	}
	temp.RightChild = node
	node.Parent = temp
	return
}

func (tree *RBtree) createNode(key Item) *Node {

	node := RBNodeAlloc()
	if node == nil {
		return nil
	}
	node.KeyValue = key
	node.Parent = &LEAF
	node.LeftChild = &LEAF
	node.RightChild = &LEAF
	/* by default every new node is red */
	node.color = RED
	return node
}

func RBNodeAlloc() *Node {
	temp := rbPool.Alloc(NODESIZE)
	node := (*Node)(Byte2Pointer(temp))
	node.addr = temp
	return node
}

func RBNodeFree(node *Node) {
	rbPool.Free(node.addr)
}

func InitRBtreeMemPool() {
	rbPool = slab.NewSyncPool(
		60,          // The smallest chunk size is 64B.
		60*1024*200, // The largest chunk size is 64KB.
		2,           // Power of 2 growth in chunk size.
	)

}
