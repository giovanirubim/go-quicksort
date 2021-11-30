package avl_tree

import "strconv"

type Key int
type Value int
type Pair struct {
	Key   Key
	Value Value
}

func CompareKeys(a, b Key) int {
	return int(a) - int(b)
}

func (a *Pair) CompareTo(b *Pair) int {
	return CompareKeys(a.Key, b.Key)
}

func (key Key) ToString() string {
	return strconv.Itoa(int(key))
}
