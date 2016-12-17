package heap

import (
	"fmt"
	"testing"
)

func print_heap(heap *Heap) {

	fmt.Printf("Heap size is %d\n", heap.Size())

	for i := uint32(0); i < heap.Size(); i++ {
		fmt.Printf("Node=%d\n", *(heap.tree[i].(*int)))
	}
	return

}

func compare_int(int1, int2 EleType) int {

	k1, _ := int1.(*int)
	k2, _ := int2.(*int)

	//fmt.Printf("%d-----%d\n", *k1, *k2)
	if *k1 < *k2 {
		fmt.Println("return 1")
		return 1
	} else if *k1 > *k2 {
		fmt.Println("return -1")
		return -1
	} else {
		fmt.Println("return 0")
		return 0
	}
}

func TestHeap(t *testing.T) {

	var intval [30]int

	heap := &Heap{}
	heap.Init(compare_int, nil, 100)

	i := 0

	intval[i] = 5
	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	intval[i] = 10
	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}

	print_heap(heap)
	i++

	intval[i] = 20
	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	intval[i] = 1
	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	intval[i] = 25
	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	intval[i] = 22

	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	intval[i] = 9

	fmt.Println("Inserting ", intval[i])
	if heap.Insert(&intval[i]) != 0 {
		t.Errorf("heap Insert fail")
	}
	print_heap(heap)
	i++

	first_data := heap.First()
	fmt.Printf("first data %d\n", *(first_data.(*int)))
	for heap.Size() > 0 {

		data := heap.Extract()
		fmt.Printf("Extracting %d\n", *(data.(*int)))
		print_heap(heap)
	}

	fmt.Printf("Destroying the heap\n")
	heap.Destroy()

}
