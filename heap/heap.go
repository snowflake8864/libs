package heap

import (
	//	"constant"
	"fmt"
)

type EleType interface{}

type CompFunc func(key1, key2 EleType) int
type Heap struct {
	size    uint32
	compare CompFunc
	//compare func(key1, key2 EleType) int
	destroy func(data EleType)
	tree    []EleType
}

func (heap *Heap) Size() uint32 {

	return heap.size
}

func (heap *Heap) First() EleType {
	if heap.Size() > 0 {
		return heap.tree[0]
	} else {
		return nil
	}
}
func (heap *Heap) Parent(npos uint32) uint32 {
	return uint32((npos - 1) / 2)
}

func (heap *Heap) Left(npos uint32) uint32 {
	return (2 * npos) + 1
}

func (heap *Heap) Right(npos uint32) uint32 {
	return (2 * npos) + 2
}

func (heap *Heap) Init(compare CompFunc,
	destroy func(data EleType), size uint32) {

	heap.size = 0
	heap.compare = compare
	heap.destroy = destroy
	heap.tree = make([]EleType, size)

}

func (heap *Heap) Destroy() {

	if heap.destroy != nil {

		for i := uint32(0); i < heap.Size(); i++ {

			heap.destroy(heap.tree[i])

		}

	}
	//	free(heap.tree)
	return

}
func (heap *Heap) Insert(data EleType) int {

	heap.tree[heap.Size()] = data

	ipos := heap.Size()
	ppos := heap.Parent(ipos)

	for ipos > 0 && heap.compare(heap.tree[ppos], heap.tree[ipos]) < 0 {

		temp := heap.tree[ppos]
		heap.tree[ppos] = heap.tree[ipos]
		heap.tree[ipos] = temp

		ipos = ppos
		ppos = heap.Parent(ipos)
		fmt.Println("-------Insert---------")
	}

	heap.size++

	return 0
}
func (heap *Heap) Peek(pos uint32) EleType {

	if heap.Size() < pos {
		return nil
	}
	return heap.tree[pos]
}
func (heap *Heap) Extract() EleType {

	if heap.Size() == 0 {
		return -1
	}

	data := heap.tree[0]

	save := heap.tree[heap.Size()-1]

	if heap.Size()-1 > 0 {

		heap.size--

	} else {
		heap.size = 0
		return data
	}
	heap.tree[0] = save

	mpos := uint32(0)
	ipos := uint32(0)
	lpos := heap.Left(ipos)
	rpos := heap.Right(ipos)

	for {

		lpos = heap.Left(ipos)
		rpos = heap.Right(ipos)

		if lpos < heap.Size() && heap.compare(heap.tree[lpos], heap.tree[ipos]) > 0 {

			mpos = lpos

		} else {

			mpos = ipos

		}
		if rpos < heap.Size() && heap.compare(heap.tree[rpos], heap.tree[mpos]) > 0 {

			mpos = rpos

		}

		if mpos == ipos {

			break

		} else {
			temp := heap.tree[mpos]
			heap.tree[mpos] = heap.tree[ipos]
			heap.tree[ipos] = temp

			ipos = mpos

		}

	}
	return data

}
