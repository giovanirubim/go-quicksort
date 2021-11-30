package avl_tree

type AVLTree struct {
	root *node
	Size int
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

func (tree *AVLTree) Set(pair Pair) {
	root, added := tree.root.Set(pair)
	tree.root = root
	if added {
		tree.Size++
	}
}

func (tree *AVLTree) Get(key Key) (Value, bool) {
	return tree.root.Get(key)
}

func (tree *AVLTree) Delete(key Key) bool {
	var deleted bool
	tree.root, deleted = tree.root.Delete(key)
	return deleted
}

func (tree *AVLTree) PairChannel() chan Pair {
	channel := make(chan Pair)
	go func() {
		tree.root.Iterate(channel)
		close(channel)
	}()
	return channel
}
