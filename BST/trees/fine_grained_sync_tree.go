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

func (fNd *FineNode[T, K]) Lock() {
	fNd.mutex.Lock()
}

func (fNd *FineNode[T, K]) Unlock() {
	fNd.mutex.Unlock()
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

func (t *FineGrainedSyncTree[T, K]) UnlockParent(parent *FineNode[T, K]) {
	if parent == nil {
		t.mutex.Unlock()
	} else {
		parent.Unlock()
	}
}

func (t *FineGrainedSyncTree[T, K]) Remove(key K) {
	currNode, parentNode := t.FinderNode(key)

	defer t.UnlockParent(parentNode)

	if currNode == nil {
		return
	}

	switch {
	case currNode.left == nil && currNode.right == nil:
		if currNode == t.root {
			t.root = nil
		} else if parentNode.left != nil && parentNode.left == currNode {
			parentNode.left = nil
		} else {
			parentNode.right = nil
		}
	case currNode.left != nil && currNode.right == nil:
		if currNode == t.root {
			t.root = currNode.left
		} else if parentNode.left != nil && parentNode.left == currNode {
			parentNode.left = currNode.left
		} else {
			parentNode.right = currNode.left
		}

	case currNode.right != nil && currNode.left == nil:
		if currNode == t.root {
			t.root = currNode.right
		} else if parentNode.left != nil && parentNode.left == currNode {
			parentNode.left = currNode.right
		} else {
			parentNode.right = currNode.right
		}
	default:
		// 2 child nodes in current Node
		defer currNode.Unlock()
		currNode.right.Lock()

		tmpParent := currNode
		tmpNode := currNode.right
		for tmpNode.left != nil {
			tmpGrandParent := tmpParent
			tmpParent = tmpNode
			tmpNode.left.Lock()
			tmpNode = tmpNode.left
			if tmpGrandParent != currNode {
				tmpGrandParent.Unlock()
			}
		}

		if tmpParent != currNode {
			defer tmpParent.Unlock()
			tmpParent.left = tmpNode.right
		} else {
			tmpParent.right = tmpNode.right
		}
		currNode.value = tmpNode.value
		currNode.key = tmpNode.key
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

func (t *FineGrainedSyncTree[T, K]) IsValid() bool {
	return t.root.isValid()
}

func (fNd *FineNode[T, K]) isValid() bool {
	if fNd == nil {
		return true
	}
	if fNd.left != nil && fNd.left.key >= fNd.key {
		return false
	}
	if fNd.right != nil && fNd.right.key <= fNd.key {
		return false
	}
	return fNd.left.isValid() && fNd.right.isValid()
}
