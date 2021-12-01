package avl_tree

type AVLTree struct {
	root *node
	size int
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

func (tree *AVLTree) Set(pair Pair) {
	root, added := tree.root.Set(pair)
	tree.root = root
	if added {
		tree.size++
	}
}

func (tree *AVLTree) Get(key Key) (Value, bool) {
	return tree.root.Get(key)
}

func (tree *AVLTree) Delete(key Key) bool {
	var deleted bool
	tree.root, deleted = tree.root.Delete(key)
	if deleted {
		tree.size--
	}
	return deleted
}

func (tree *AVLTree) Size() int {
	return tree.size
}

func (tree *AVLTree) PairChannel() chan Pair {
	channel := make(chan Pair)
	go func() {
		tree.root.Iterate(channel)
		close(channel)
	}()
	return channel
}
