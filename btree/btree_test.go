package btree

import (
	"fmt"
	"testing"
)

func TestBtree(t *testing.T) {
	bt := NewBTree(3)
	for i := range 5 {
		bt.Put(&Data{i + 1, i, false})
		bt.showDetail()
	}
	bt.PrintTree()
	bt.Put(&Data{6, 6, false})
	bt.showDetail()
	bt.Put(&Data{7, 7, false})
	bt.showDetail()
	bt.Remove(&Data{4, 7, false})
	bt.showDetail()
	ans := bt.Get(&Data{4, 7, false})
	fmt.Println(ans)
	bt.showDetail()
	bt.PrintTree()
}
