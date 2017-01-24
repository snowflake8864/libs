package public

import (
	"fmt"
	"testing"
)

func testfn(a, b int) (int, Derror) {
	if a > b {
		return 10, nil
	} else {

		return -1, Perror(ERR_EPERM)
	}
}

func TestPublic(t *testing.T) {
	fmt.Println(testfn(50, 30))
}
