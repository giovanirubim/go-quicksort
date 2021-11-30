package avl_tree

import (
	"testing"
)

func ToPair(value int) Pair {
	return Pair{Key(value), Value(value)}
}

func VerifySorting(node *node, t *testing.T) bool {
	if node == nil {
		return true
	}

	left, right := node.Left, node.Right

	if left != nil && left.Pair.CompareTo(&node.Pair) > 0 {
		t.Error("Node", node.Pair, "has", node.Left.Pair, "on its left")
		return false
	}
	if right != nil && right.Pair.CompareTo(&node.Pair) < 0 {
		t.Error("Node", node.Pair, "has", node.Right.Pair, "on its right")
		return false
	}
	if !VerifySorting(left, t) {
		return false
	}
	if !VerifySorting(right, t) {
		return false
	}
	return true
}

func VerifyHeightsRec(node *node, t *testing.T) (bool, int) {
	if node == nil {
		return true, -1
	}

	lMatches, lHeight := VerifyHeightsRec(node.Left, t)
	rMatches, rHeight := VerifyHeightsRec(node.Right, t)
	height := max(lHeight, rHeight) + 1

	if node.Height != height {
		t.Error("Wrong height at", node.Pair, "\nExpected:", height, "actual:", node.Height)
		return false, height
	}
	return lMatches && rMatches, height
}

func VerifyHeights(node *node, t *testing.T) bool {
	ok, _ := VerifyHeightsRec(node, t)
	return ok
}

func VerifyBalancing(node *node, t *testing.T) bool {
	if node == nil {
		return true
	}

	left, right := -1, -1
	if node.Left != nil {
		left = node.Left.Height
	}
	if node.Right != nil {
		right = node.Right.Height
	}

	factor := left - right
	if factor < -1 {
		t.Error("Tree two heavy twords the right at", node.Pair)
		return false
	}
	if factor > 1 {
		t.Error("Tree two heavy twords the left at", node.Pair)
		return false
	}
	return true
}

func VerifyAll(root *node, t *testing.T) bool {
	heightOk := VerifyHeights(root, t)
	sorted := VerifySorting(root, t)
	balanced := VerifyBalancing(root, t)
	return heightOk && sorted && balanced
}

func VerifyString(root *node, expected string, t *testing.T) {
	actual := root.ToString()
	if actual != expected {
		t.Error("Wrong string\nExpected: " + expected + "\nActual:   " + actual)
	}
}

func VerifyInsertions(values []int, t *testing.T) (*node, bool) {
	var root *node = nil
	for i, value := range values {
		root, _ = root.Set(ToPair(value))
		if !VerifyAll(root, t) {
			t.Error("Failed after inserting", values[:i+1])
			return nil, false
		}
	}
	return root, true
}

func Populate(values []int) *node {
	var root *node = nil
	for _, value := range values {
		root, _ = root.Set(ToPair(value))
	}
	return root
}

func CalcSize(root *node) int {
	if root == nil {
		return 0
	}
	return CalcSize(root.Left) + 1 + CalcSize(root.Right)
}

func VerifySize(root *node, expected int, t *testing.T) bool {
	actual := CalcSize(root)

	if actual != expected {
		t.Error("Wrong size\nExpected:", expected, "\nActual:", actual)
		return false
	}

	return true
}

func TestNewNode(t *testing.T) {
	node := newNode(Pair{1, 2})

	if node.Height != 0 {
		t.Error("Height of new node is not 0")
	}
	if node.Pair != (Pair{1, 2}) {
		t.Error("Did not copy pair values correctly")
	}
	if node.Left != nil || node.Right != nil {
		t.Error("New node is not a leaf")
	}
}

func TestAddingValueToNilNode(t *testing.T) {
	var root *node = nil
	var added bool
	root, added = root.Set(ToPair(1))

	if root == nil {
		t.Error("Adding element returned nil")
	}
	if !added {
		t.Error("Did not return added as true")
	}
	if root.Pair != (ToPair(1)) {
		t.Error("Failed to insert correct values")
	}
	if root.Left != nil || root.Right != nil {
		t.Error("Inserted node is not a leaf")
	}
}

func TestAddingHigherValue(t *testing.T) {
	var root *node = nil
	root, _ = root.Set(ToPair(5))
	root, _ = root.Set(ToPair(9))

	if root.Height != 1 {
		t.Errorf("Wrong tree height\nExpected: %d\nActual: %d", 1, root.Height)
	}
	if root.Pair != ToPair(5) {
		t.Error("Wrong root element")
	}
	if root.Right == nil {
		t.Error("Did not insert new element to the right")
	}
	if root.Left != nil {
		t.Error("Inserted something to the left")
	}
	if root.Right.Pair != ToPair(9) {
		t.Error("Did not insert new pair correctly")
	}
}

func TestAddingLowerValue(t *testing.T) {
	var root *node = nil
	root, _ = root.Set(ToPair(5))
	root, _ = root.Set(ToPair(1))

	if root.Height != 1 {
		t.Errorf("Wrong tree height\nExpected: %d\nActual: %d", 1, root.Height)
	}
	if root.Pair != ToPair(5) {
		t.Error("Wrong root element")
	}
	if root.Left == nil {
		t.Error("Did not insert new element to the left")
	}
	if root.Right != nil {
		t.Error("Inserted something to the right")
	}
	if root.Left.Pair != ToPair(1) {
		t.Error("Did not insert new pair correctly")
	}
	VerifySize(root, 2, t)
}

func TestSimpleLeftRotation(t *testing.T) {
	var root *node = nil
	root, _ = root.Set(ToPair(2))
	root, _ = root.Set(ToPair(1))
	root, _ = root.Set(ToPair(4))
	root, _ = root.Set(ToPair(3))
	root, _ = root.Set(ToPair(5))

	root = root.RotateLeft()

	if root.Height != 2 {
		t.Errorf("Wrong height\nExpected: 2\nActual %d\n", root.Height)
	}
	if root.Pair.Key != 4 {
		t.Errorf("Wrong root key\nExpected: 4\nActual %d\n", root.Pair.Key)
	}

	factor := root.Left.Height - root.Right.Height
	if factor != 1 {
		t.Errorf("Wrong balancing factor\nExpected 1\nActual %d\n", factor)
	}

	VerifySorting(root, t)
	VerifyHeights(root, t)
	VerifySize(root, 5, t)
}

func TestSimpleRightRotation(t *testing.T) {
	var root *node = nil
	root, _ = root.Set(ToPair(4))
	root, _ = root.Set(ToPair(2))
	root, _ = root.Set(ToPair(5))
	root, _ = root.Set(ToPair(1))
	root, _ = root.Set(ToPair(3))

	root = root.RotateRight()

	if root.Height != 2 {
		t.Errorf("Wrong height\nExpected: 2\nActual %d\n", root.Height)
	}
	if root.Pair.Key != 2 {
		t.Errorf("Wrong root key\nExpected: 2\nActual %d\n", root.Pair.Key)
	}

	factor := root.Left.Height - root.Right.Height
	if factor != -1 {
		t.Errorf("Wrong balancing factor\nExpected -1\nActual %d\n", factor)
	}

	VerifySorting(root, t)
	VerifyHeights(root, t)
	VerifySize(root, 5, t)
}

func TestSimpleBalancing(t *testing.T) {
	tests := [][]int{
		{1, 2, 3},
		{3, 2, 1},
	}
	for _, values := range tests {
		root, ok := VerifyInsertions(values, t)
		if ok && !VerifySize(root, len(values), t) {
			t.Error("Failed after inserting", values)
		}
	}
}

func TestDoubleRotationBalancing(t *testing.T) {
	tests := [][]int{
		{1, 3, 2},
		{3, 1, 2},
		{2, 1, 3, 4, 6, 5},
		{2, 1, 5, 4, 6, 3},
	}
	for _, values := range tests {
		root, ok := VerifyInsertions(values, t)
		if ok && !VerifySize(root, len(values), t) {
			t.Error("Failed after inserting", values)
		}
	}
}

func TestGetPresentElement(t *testing.T) {
	root, ok := VerifyInsertions([]int{2, 1, 5, 4, 6, 3}, t)
	if !ok {
		return
	}

	value, found := root.Get(Key(3))
	if !found {
		t.Error("Key", 3, "not found in tree")
		return
	}

	if value != 3 {
		t.Error("Returned value differs from inserted one")
		return
	}
}

func TestUpdate(t *testing.T) {
	root, ok := VerifyInsertions([]int{1, 2, 3, 4, 5, 6, 7}, t)
	if !ok {
		return
	}

	targetKey := Key(5)
	newValue := Value(9)

	root.Set(Pair{targetKey, newValue})
	VerifyAll(root, t)
	VerifySize(root, 7, t)

	value, found := root.Get(Key(targetKey))

	if !found {
		t.Error("Updated item not found")
		return
	}

	if value != newValue {
		t.Error("Did not update the pair correctly")
		return
	}
}

func TestRemoveLeftest(t *testing.T) {
	root, ok := VerifyInsertions([]int{3, 1, 6, 2, 5, 7, 4}, t)
	if !ok {
		return
	}

	var pair Pair
	root, pair = root.RemoveLeftest()

	if pair != ToPair(1) {
		t.Error("Removed incorrect pair\nExpected:", ToPair(1), "\nActual:", pair)
	}

	VerifyAll(root, t)
	VerifySize(root, 6, t)
}

func TestToString(t *testing.T) {
	root, ok := VerifyInsertions([]int{1, 2, 3, 4, 5, 6, 7}, t)
	if !ok {
		return
	}

	VerifyString(root, "(((-,1,-),2,(-,3,-)),4,((-,5,-),6,(-,7,-)))", t)
}

func TestLeafDeletion(t *testing.T) {
	root := Populate([]int{3, 1, 4, 2})
	target := Key(2)

	var deleted bool
	root, deleted = root.Delete(target)

	if !deleted {
		t.Error("Returned deleted as false")
	}

	VerifyAll(root, t)
	VerifySize(root, 3, t)

	_, found := root.Get(target)
	if found {
		t.Error("Returned found after deletion")
	}
}

func TestLeafDeletionWithRebalancing(t *testing.T) {
	root := Populate([]int{3, 1, 4, 2})
	target := Key(4)

	var deleted bool
	root, deleted = root.Delete(target)

	if !deleted {
		t.Error("Returned deleted as false")
	}

	VerifyAll(root, t)
	VerifySize(root, 3, t)

	_, found := root.Get(target)
	if found {
		t.Error("Returned found after deletion")
	}
}

func TestRootDeletion(t *testing.T) {
	root := Populate([]int{2, 1, 3, 4})
	target := Key(2)

	var deleted bool
	root, deleted = root.Delete(target)

	if !deleted {
		t.Error("Returned deleted as false")
	}

	VerifyAll(root, t)
	VerifySize(root, 3, t)

	_, found := root.Get(target)
	if found {
		t.Error("Returned found after deletion")
	}
}

func TestRootDeletionWithRebalancing(t *testing.T) {
	root := Populate([]int{3, 1, 4, 2})
	target := Key(3)

	var deleted bool
	root, deleted = root.Delete(target)

	if !deleted {
		t.Error("Returned deleted as false")
	}

	VerifyAll(root, t)
	VerifySize(root, 3, t)

	_, found := root.Get(target)
	if found {
		t.Error("Returned found after deletion")
	}
}

func TestIterate(t *testing.T) {
	root := Populate([]int{3, 1, 4, 1, 5, 9, 2, 6, 5})
	VerifyAll(root, t)

	channel := make(chan Pair)
	go func() {
		defer close(channel)
		root.Iterate(channel)
	}()

	expected := []Pair{
		ToPair(1),
		ToPair(2),
		ToPair(3),
		ToPair(4),
		ToPair(5),
		ToPair(6),
		ToPair(9),
	}

	actual := []Pair{}
	for pair := range channel {
		actual = append(actual, pair)
	}

	if len(actual) != len(expected) {
		t.Error("Return wrong amount of pairs")
	}
	for i, value := range expected {
		if actual[i] != value {
			t.Error("Returned wrong sequence of pairs\nExpected:", expected, "\nActual:  ", actual)
			break
		}
	}
}
