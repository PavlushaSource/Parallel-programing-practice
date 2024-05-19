package trees

import (
	"cmp"
	"sync"
)

type GrainedSyncTree[T any, K cmp.Ordered] struct {
	root  *Node[T, K]
	mutex *sync.Mutex
}

type Node[T any, K cmp.Ordered] struct {
	key   K
	value T
	left  *Node[T, K]
	right *Node[T, K]
}

func NewGrainedSyncTree[T any, K cmp.Ordered]() *GrainedSyncTree[T, K] {
	return &GrainedSyncTree[T, K]{
		root:  nil,
		mutex: &sync.Mutex{},
	}
}
func (t *GrainedSyncTree[T, K]) find(key K) (value T, exist bool) {
	node := t.root

	for node != nil {
		switch cmp.Compare(key, node.key) {
		case -1:
			node = node.left
		case 1:
			node = node.right
		case 0:
			return node.value, true
		}
	}
	return value, false
}

func (t *GrainedSyncTree[T, K]) Find(key K) (T, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.find(key)
}

func (t *GrainedSyncTree[T, K]) Insert(key K, value T) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.root = insert(t.root, key, value)
}

func insert[T any, K cmp.Ordered](node *Node[T, K], key K, value T) *Node[T, K] {
	if node == nil {
		return &Node[T, K]{key: key, value: value}
	}
	switch cmp.Compare(key, node.key) {
	case -1:
		node.left = insert(node.left, key, value)
	case 1:
		node.right = insert(node.right, key, value)
	case 0:
		node.value = value
	}
	return node
}

func (t *GrainedSyncTree[T, K]) Remove(key K) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.root = t.remove(key, t.root)
}

func (t *GrainedSyncTree[T, K]) remove(key K, node *Node[T, K]) *Node[T, K] {

	if node == nil {
		return nil
	}

	switch cmp.Compare(key, node.key) {
	case -1:
		node.left = t.remove(key, node.left)
	case 1:
		node.right = t.remove(key, node.right)
	case 0:
		if node.left != nil && node.right != nil {
			minRightNode := t.min(node.right)
			node.key = minRightNode.key
			node.value = minRightNode.value
			node.right = t.remove(minRightNode.key, node.right)
		} else if node.left == nil {
			node = node.right
		} else {
			node = node.left
		}
	}

	return node
}

func (t *GrainedSyncTree[T, K]) min(node *Node[T, K]) *Node[T, K] {
	for node.left != nil {
		node = node.left
	}
	return node
}
