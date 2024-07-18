package btree

import (
	"fmt"
)

// 注意 b+树的实现参考的是 可视化的逻辑进行的
// https://www.cs.usfca.edu/~galles/visualization/BPlusTree.html
// 但这个和我们平时印象中的b+树有点差距的
// 但可以根据具体的应用场景进行切换 这样实现的话 可以最大程度避免 最大节点向上传递的冗余
// 对于删除操作 后续进行补充实现
// 目前使用懒删除 标记删除即可
// 已经达到学习要求了

type BPTreeNode struct {
	entries  DataItem
	children []*BPTreeNode
	parent   *BPTreeNode
}
type BPTree struct {
	root  *BPTreeNode
	order int
	size  int
}

// 这部分存在差异 但并不影响
func (bt *BPTree) max_entries() int {
	return bt.order - 1
}

func (bt *BPTree) middle() int {
	return (bt.order - 1) / 2
}
func (bt *BPTree) isLeaf(node *BPTreeNode) bool {
	return len(node.children) == 0
}
func (bt *BPTree) shouldSplit(node *BPTreeNode) bool {
	return len(node.entries) > bt.max_entries()
}

func (bt *BPTree) levelOrderTraversal() (ans []DataItem) {
	root := bt.root
	if root == nil {
		return
	}
	queue := []*BPTreeNode{root}
	for i := 0; i < len(queue); i++ {
		ans = append(ans, []*Data{})
		var pueue []*BPTreeNode
		for j := 0; j < len(queue); j++ {
			node := queue[j]
			ans[i] = append(ans[i], node.entries...)
			pueue = append(pueue, node.children...)
		}
		queue = pueue
	}
	return
}
func (bt *BPTree) Put(item *Data) {
	if bt.root == nil {
		bt.root = &BPTreeNode{entries: DataItem{item}}
		return
	}
	if bt.insert(bt.root, item) {
		bt.size++
	}
}
func (bt *BPTree) delete(node *BPTreeNode, idx int) {
	node.entries[idx].deleted = true
	bt.size--
}
func (bt *BPTree) Remove(item *Data) bool {
	node, idx, ok := bt.search(bt.root, item)
	if ok {
		bt.delete(node, idx)
	}
	return ok
}
func (bt *BPTree) search(node *BPTreeNode, item *Data) (st *BPTreeNode, idx int, ok bool) {
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
func (bt *BPTree) Get(item *Data) bool {
	node, index, ok := bt.search(bt.root, item)
	return ok && !node.entries[index].deleted
}
func (bt *BPTree) insertLeaf(node *BPTreeNode, item *Data) bool {
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
func (bt *BPTree) split(node *BPTreeNode) {
	if !bt.shouldSplit(node) {
		return
	}
	if bt.root == node {
		bt.splitRoot()
		return
	}
	bt.splitNonRoot(node)
}

func (bt *BPTree) splitRoot() {
	mid := bt.middle()
	node := bt.root

	midItem := node.entries[mid]
	newRoot := &BPTreeNode{entries: DataItem{midItem}}

	leftNode := &BPTreeNode{entries: node.entries[:mid], parent: newRoot}
	// rightNode := &BPTreeNode{entries: node.entries[mid+1:], parent: newRoot}
	var rightNode *BPTreeNode
	if bt.isLeaf(node) {
		rightNode = &BPTreeNode{entries: node.entries[mid:], parent: newRoot}
	} else {
		rightNode = &BPTreeNode{entries: node.entries[mid+1:], parent: newRoot}
	}

	if !bt.isLeaf(node) {
		leftNode.children = node.children[:mid+1]
		rightNode.children = node.children[mid+1:]
		setParentBPTree(leftNode.children, leftNode)
		setParentBPTree(rightNode.children, rightNode)
	}
	newRoot.children = []*BPTreeNode{leftNode, rightNode}
	bt.root = newRoot
}

func (bt *BPTree) splitNonRoot(node *BPTreeNode) {
	mid := bt.middle()
	parent := node.parent

	midItem := node.entries[mid]

	leftNode := &BPTreeNode{entries: node.entries[:mid], parent: parent}
	// rightNode := &BPTreeNode{entries: node.entries[mid+1:], parent: parent}
	rightNode := &BPTreeNode{entries: node.entries[mid:], parent: parent}
	if !bt.isLeaf(node) {
		leftNode.children = node.children[:mid+1]
		rightNode.children = node.children[mid+1:]
		setParentBPTree(leftNode.children, leftNode)
		setParentBPTree(rightNode.children, rightNode)
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
func (bt *BPTree) insertInnner(node *BPTreeNode, item *Data) bool {
	index, ok := node.entries.search(item)
	if ok {
		node.entries.changeData(index, item)
		return false
	}
	return bt.insert(node.children[index], item)
}
func (bt *BPTree) insert(node *BPTreeNode, item *Data) bool {
	if bt.isLeaf(node) {
		return bt.insertLeaf(node, item)
	}
	return bt.insertInnner(node, item)
}
func NewBPTree(order int) *BPTree {
	return &BPTree{order: max(order, 3)}
}
func (bt *BPTree) showDetail() {
	fmt.Println()
	ans := bt.levelOrderTraversal()
	for _, v := range ans {
		fmt.Println("[" + v.showInfo() + "]")
	}
	fmt.Println()
}
func (bt *BPTree) PrintTree() {
	bt.printTree(bt.root, 0, " ")
}
func (bt *BPTree) printTree(node *BPTreeNode, depth int, prefix string) {
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
func childPrefix(isLastChild bool) string {
	if isLastChild {
		return "└-- "
	}
	return "|-- "
}
