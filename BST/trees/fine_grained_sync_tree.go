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
	currNode, parentNode := t.FinderNode(key)
	insertNode := &FineNode[T, K]{key: key, value: value, mutex: &sync.Mutex{}}

	if parentNode == nil {
		if currNode != nil {
			currNode.value = value
			defer currNode.Unlock()
		} else {
			t.root = insertNode
		}
		t.mutex.Unlock()
		return

	} else {
		defer parentNode.Unlock()
		if currNode != nil {
			currNode.value = value
			currNode.Unlock()
			return
		} else {
			switch cmp.Compare(key, parentNode.key) {
			case -1:
				parentNode.left = insertNode
			case 1:
				parentNode.right = insertNode
			default:
				panic("this should not happen: parent.key = insert key")
			}
			return
		}
	}
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
		value = currNode.value
		return value, true
	}
	return
}

func (t *FineGrainedSyncTree[T, K]) Remove(key K) {
	currNode, parentNode := t.FinderNode(key)

	if currNode == nil {
		if parentNode != nil {
			parentNode.Unlock()
		} else {
			t.mutex.Unlock()
		}
		return
	} else {
		switch {
		case currNode.left == nil && currNode.right == nil:
			if parentNode.left != nil && parentNode.left == currNode {
				parentNode.left = nil
			} else {
				parentNode.right = nil
			}
			parentNode.Unlock()
			return
		case currNode.left != nil && currNode.right == nil:
			if parentNode.left != nil && parentNode.left == currNode {
				parentNode.left = currNode.left
			} else {
				parentNode.right = currNode.left
			}
			parentNode.Unlock()
			return
		case currNode.right != nil && currNode.left == nil:
			if parentNode.left != nil && parentNode.left == currNode {
				parentNode.left = currNode.right
			} else {
				parentNode.right = currNode.right
			}
			parentNode.Unlock()
			return
		default:
			// 2 child nodes current Node
			currNode.right.Lock()
			currNode.left.Lock()

			rChild := currNode.right
			lChild := currNode.left
			if currNode.left.right == nil {
				currNode.left.right = rChild
			} else {
				subTree := FineGrainedSyncTree[T, K]{mutex: &sync.Mutex{}, root: currNode.left.right}
				subTree.Insert(rChild.key, rChild.value)
			}

		}
	}

}

func NewFineNode[T any, K cmp.Ordered]() *FineNode[T, K] {
	return &FineNode[T, K]{mutex: &sync.Mutex{}}
}

func (t *FineGrainedSyncTree[T, K]) FinderNode(key K) (currentNode *FineNode[T, K], parentNode *FineNode[T, K]) {
	t.mutex.Lock()

	if t.root == nil {
		return nil, nil
	}

	t.root.Lock()
	currentNode = t.root

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
			return parentNode, grandParent
		}

		// Анлочим древо
		if grandParent == nil {
			t.mutex.Unlock()
		}

		// Анлочим деда
		if grandParent != nil {
			grandParent.Unlock()
		}

	}
	return
}
