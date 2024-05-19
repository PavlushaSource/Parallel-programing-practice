package trees

import (
	"cmp"
	"sync"
)

type FineGrainedSyncTree[T any, K cmp.Ordered] struct {
	root  *FineNode[T, K]
	mutex *sync.Mutex
}

type FineNode[T any, K cmp.Ordered] struct {
	key   K
	value T
	left  *FineNode[T, K]
	right *FineNode[T, K]
	mutex *sync.Mutex
}

func (fn *FineNode[T, K]) Lock() {
	fn.mutex.Lock()
}

func (fn *FineNode[T, K]) Unlock() {
	fn.mutex.Unlock()
}

func NewFineGrainedSyncTree[T any, K cmp.Ordered]() *FineGrainedSyncTree[T, K] {
	return &FineGrainedSyncTree[T, K]{
		root:  nil,
		mutex: &sync.Mutex{},
	}
}

func (t *FineGrainedSyncTree[T, K]) Insert(key K, value T) {
	// TODO IT
}

func (t *FineGrainedSyncTree[T, K]) Find(key K) (value T, exist bool) {
	currNode, parentNode := t.FinderNode(key)

	if parentNode == nil {
		defer t.mutex.Unlock()
	} else {
		defer parentNode.Unlock()
	}

	if currNode != nil {
		defer currNode.Unlock()
		return currNode.value, true
	}
	return
}

func (t *FineGrainedSyncTree[T, K]) Remove(key K) {
	// TODO IT
}

func (t *FineGrainedSyncTree[T, K]) FinderNode(key K) (currentNode *FineNode[T, K], parentNode *FineNode[T, K]) {
	t.mutex.Lock()

	if t.root == nil {
		return nil, nil
	}

	t.root.Lock()
	currentNode = t.root
	parentNode = new(FineNode[T, K])

	for currentNode != nil {
		grandParent := parentNode
		parentNode = currentNode

		switch cmp.Compare(key, currentNode.key) {
		case -1:
			if currentNode.left != nil {
				currentNode.left.Lock()
			}
			currentNode = currentNode.left
		case 1:
			if currentNode.right != nil {
				currentNode.right.Lock()
			}
			currentNode = currentNode.right
		case 0:
			return
		}

		// Анлочим древо
		if grandParent == nil {
			t.root.mutex.Unlock()
		}

		// Анлочим деда
		if grandParent != nil {
			grandParent.Unlock()
		}

	}
	return
}
