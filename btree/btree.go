package btree

import (
	"fmt"
)

// 注意 b树的实现参考的是 可视化的逻辑进行的
// https://www.cs.usfca.edu/~galles/visualization/BTree.html
// 但可以根据具体的应用场景进行切换
// 对于删除操作 后续进行补充实现
// 目前使用懒删除 标记删除即可
// 已经达到学习要求了

type BTreeNode struct {
	entries  DataItem
	children []*BTreeNode
	parent   *BTreeNode
}
type BTree struct {
	root  *BTreeNode
	order int
	size  int
}

func (bt *BTree) max_entries() int {
	return bt.order - 1
}

func (bt *BTree) middle() int {
	return (bt.order - 1) / 2
}
func (bt *BTree) isLeaf(node *BTreeNode) bool {
	return len(node.children) == 0
}
func (bt *BTree) shouldSplit(node *BTreeNode) bool {
	return len(node.entries) > bt.max_entries()
}

func (bt *BTree) levelOrderTraversal() (ans []DataItem) {
	root := bt.root
	if root == nil {
		return
	}
	queue := []*BTreeNode{root}
	for i := 0; i < len(queue); i++ {
		ans = append(ans, []*Data{})
		var pueue []*BTreeNode
		for j := 0; j < len(queue); j++ {
			node := queue[j]
			ans[i] = append(ans[i], node.entries...)
			pueue = append(pueue, node.children...)
		}
		queue = pueue
	}
	return
}
func (bt *BTree) Put(item *Data) {
	if bt.root == nil {
		bt.root = &BTreeNode{entries: DataItem{item}}
		return
	}
	if bt.insert(bt.root, item) {
		bt.size++
	}
}
func (bt *BTree) delete(node *BTreeNode, idx int) {
	node.entries[idx].deleted = true
	bt.size--
}
func (bt *BTree) Remove(item *Data) bool {
	node, idx, ok := bt.search(bt.root, item)
	if ok {
		bt.delete(node, idx)
	}
	return ok
}
func (bt *BTree) search(node *BTreeNode, item *Data) (st *BTreeNode, idx int, ok bool) {
	if bt.root == nil {
		return nil, -1, false
	}
	st = node
	for {
		idx, ok = node.entries.search(item)
		if ok {
			return node, idx, true
		}
		if bt.isLeaf(node) {
			return nil, -1, false
		}
		node = node.children[idx]
	}
}
func (bt *BTree) Get(item *Data) bool {
	node, index, ok := bt.search(bt.root, item)
	return ok && !node.entries[index].deleted
}
func (bt *BTree) insertLeaf(node *BTreeNode, item *Data) bool {
	index, ok := node.entries.search(item)
	if ok {
		// node.entries.changeData(index, item)
		node.entries[index] = item
		return false
	}
	node.entries = append(node.entries, nil)
	copy(node.entries[index+1:], node.entries[index:])
	node.entries[index] = item
	bt.split(node)
	return true
}
func (bt *BTree) split(node *BTreeNode) {
	if !bt.shouldSplit(node) {
		return
	}
	if bt.root == node {
		bt.splitRoot()
		return
	}
	bt.splitNonRoot(node)
}

func (bt *BTree) splitRoot() {
	mid := bt.middle()
	node := bt.root

	midItem := node.entries[mid]
	newRoot := &BTreeNode{entries: DataItem{midItem}}

	leftNode := &BTreeNode{entries: node.entries[:mid], parent: newRoot}
	rightNode := &BTreeNode{entries: node.entries[mid+1:], parent: newRoot}
	if !bt.isLeaf(node) {
		leftNode.children = node.children[:mid+1]
		rightNode.children = node.children[mid+1:]
		setParentBTree(leftNode.children, leftNode)
		setParentBTree(rightNode.children, rightNode)
	}
	newRoot.children = []*BTreeNode{leftNode, rightNode}
	bt.root = newRoot
}

func (bt *BTree) splitNonRoot(node *BTreeNode) {
	mid := bt.middle()
	parent := node.parent

	midItem := node.entries[mid]

	leftNode := &BTreeNode{entries: node.entries[:mid], parent: parent}
	rightNode := &BTreeNode{entries: node.entries[mid+1:], parent: parent}
	if !bt.isLeaf(node) {
		leftNode.children = node.children[:mid+1]
		rightNode.children = node.children[mid+1:]
		setParentBTree(leftNode.children, leftNode)
		setParentBTree(rightNode.children, rightNode)
	}

	index, _ := parent.entries.search(midItem)

	parent.entries = append(parent.entries, nil)
	copy(parent.entries[index+1:], parent.entries[index:])
	parent.entries[index] = midItem

	parent.children[index] = leftNode
	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = rightNode

	bt.split(parent)
}
func (bt *BTree) insertInnner(node *BTreeNode, item *Data) bool {
	index, ok := node.entries.search(item)
	if ok {
		node.entries.changeData(index, item)
		return false
	}
	return bt.insert(node.children[index], item)
}
func (bt *BTree) insert(node *BTreeNode, item *Data) bool {
	if bt.isLeaf(node) {
		return bt.insertLeaf(node, item)
	}
	return bt.insertInnner(node, item)
}
func NewBTree(order int) *BTree {
	return &BTree{order: max(order, 3)}
}
func (bt *BTree) PrintTree() {
	bt.printTree(bt.root, 0, " ")
}
func (bt *BTree) printTree(node *BTreeNode, depth int, prefix string) {
	if depth == 0 {
		fmt.Printf("+--%s\n", node.entries.showInfo())
		depth++
	}
	childCount := len(node.children)
	for idx, child := range node.children {
		last := idx == childCount-1
		newPrefix := prefix
		if last {
			newPrefix += "   "
		} else {
			newPrefix += " |  "
		}
		fmt.Println(prefix, childPrefix(last), child.entries.showInfo())
		bt.printTree(child, depth+1, newPrefix)
	}
}

func (bt *BTree) showDetail() {
	fmt.Println()
	ans := bt.levelOrderTraversal()
	for _, v := range ans {
		fmt.Println(v.showInfo())
	}
	fmt.Println()
}
