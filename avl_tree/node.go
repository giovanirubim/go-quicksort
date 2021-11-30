package avl_tree

type node struct {
	Pair        Pair
	Left, Right *node
	Height      int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newNode(pair Pair) *node {
	return &node{Pair: pair}
}

func (node *node) updateHeight() *node {
	left, right := -1, -1
	if node.Left != nil {
		left = node.Left.Height
	}
	if node.Right != nil {
		right = node.Right.Height
	}
	node.Height = max(left, right) + 1
	return node
}

func (node *node) setLeft(subtree *node) *node {
	node.Left = subtree
	return node.updateHeight()
}

func (node *node) setRight(subtree *node) *node {
	node.Right = subtree
	return node.updateHeight()
}

func (node *node) setToLeft(pair Pair) (*node, bool) {
	root, added := node.Left.Set(pair)
	return node.setLeft(root), added
}

func (node *node) setToRight(pair Pair) (*node, bool) {
	root, added := node.Right.Set(pair)
	return node.setRight(root), added
}

func (node *node) Set(pair Pair) (*node, bool) {
	if node == nil {
		return newNode(pair), true
	}
	dif := pair.CompareTo(&node.Pair)
	if dif < 0 {
		root, added := node.setToLeft(pair)
		return root.balance(), added
	}
	if dif > 0 {
		root, added := node.setToRight(pair)
		return root.balance(), added
	}
	node.Pair = pair
	return node, false
}

func (node *node) RotateLeft() *node {
	root := node.Right
	node.setRight(root.Left)
	root.setLeft(node)
	return root
}

func (node *node) RotateRight() *node {
	root := node.Left
	node.setLeft(root.Right)
	root.setRight(node)
	return root
}

func (node *node) balancingFactor() int {
	left, right := -1, -1
	if node.Left != nil {
		left = node.Left.Height
	}
	if node.Right != nil {
		right = node.Right.Height
	}
	return left - right
}

func (node *node) balance() *node {
	if node.Height < 2 {
		return node
	}

	factor := node.balancingFactor()

	if factor > 1 {
		if node.Left.balancingFactor() < 0 {
			node.setLeft(node.Left.RotateLeft())
		}
		return node.RotateRight()
	}

	if factor < -1 {
		if node.Right.balancingFactor() > 0 {
			node.setRight(node.Right.RotateRight())
		}
		return node.RotateLeft()
	}

	return node
}

func (node *node) Get(key Key) (Value, bool) {
	if node == nil {
		return Value(0), false
	}
	dif := CompareKeys(key, node.Pair.Key)
	if dif < 0 {
		return node.Left.Get(key)
	}
	if dif > 0 {
		return node.Right.Get(key)
	}
	return node.Pair.Value, true
}

func (node *node) RemoveLeftest() (*node, Pair) {
	if node.Left == nil {
		return node.Right, node.Pair
	}
	left, pair := node.Left.RemoveLeftest()
	node.setLeft(left)
	return node.balance(), pair
}

func (node *node) Delete(key Key) (*node, bool) {
	if node == nil {
		return nil, false
	}
	dif := CompareKeys(key, node.Pair.Key)
	if dif < 0 {
		left, removed := node.Left.Delete(key)
		node.setLeft(left)
		if removed {
			node = node.balance()
		}
		return node, removed
	}
	if dif > 0 {
		right, removed := node.Right.Delete(key)
		node.setRight(right)
		if removed {
			node = node.balance()
		}
		return node, removed
	}
	if node.Left == nil {
		return node.Right, true
	}
	if node.Right == nil {
		return node.Left, true
	}

	newRight, pair := node.Right.RemoveLeftest()
	node.setRight(newRight)
	node.Pair = pair
	return node.balance(), true
}

func (node *node) ToString() string {
	if node == nil {
		return "-"
	}
	res := "(" + node.Left.ToString() + ","
	res += node.Pair.Key.ToString() + ","
	res += node.Right.ToString() + ")"
	return res
}

func (node *node) Iterate(channel chan Pair) {
	if node == nil {
		return
	}
	node.Left.Iterate(channel)
	channel <- node.Pair
	node.Right.Iterate(channel)
}
