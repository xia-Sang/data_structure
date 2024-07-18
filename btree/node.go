package btree

import "fmt"

type Data struct {
	key     int
	value   interface{}
	deleted bool
}

func NewData(key int, value interface{}) *Data {
	return &Data{key, value, false}
}
func (i *Data) info() string {
	return fmt.Sprintf("(%d:%v:%v)", i.key, i.value, i.deleted)
}

type DataItem []*Data

func (dt *DataItem) showInfo() (ans string) {
	for k, v := range *dt {
		ans += fmt.Sprintf("idx:%d:%v", k, v.info()) + " "
	}
	return
}

func (dt DataItem) changeData(index int, item *Data) {
	dt[index].value = item.value
}

func (dt DataItem) search(data *Data) (int, bool) {
	left := 0
	right := len(dt) - 1
	ans := -1

	for left <= right {
		mid := (left + right) / 2
		if dt[mid].key == data.key {
			return mid, true
		} else if dt[mid].key < data.key {
			ans = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ans + 1, false
}
func setParentBPTree(nodes []*BPTreeNode, parent *BPTreeNode) {
	for _, v := range nodes {
		v.parent = parent
	}
}
func setParentBTree(nodes []*BTreeNode, parent *BTreeNode) {
	for _, v := range nodes {
		v.parent = parent
	}
}
