package trees

import (
	"cmp"
	"sync"
)

type OptimisticNode[T any, K cmp.Ordered] struct {
	key   K
	value T
	left  *OptimisticNode[T, K]
	right *OptimisticNode[T, K]
	mutex *sync.Mutex
}

type OptimisticTree[T any, K cmp.Ordered] struct {
	root  *OptimisticNode[T, K]
	mutex *sync.Mutex
}

func NewOptimisticSyncTree[T any, K cmp.Ordered]() *OptimisticTree[T, K] {
	return &OptimisticTree[T, K]{
		root:  nil,
		mutex: &sync.Mutex{},
	}
}

func (t *OptimisticTree[T, K]) Insert(key K, value T) {
	currNode, parentNode := t.FinderNode(key)
	insertNode := &OptimisticNode[T, K]{key: key, value: value, mutex: &sync.Mutex{}}

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

func (t *OptimisticTree[T, K]) Find(key K) (value T, exist bool) {
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

func (t *OptimisticTree[T, K]) Remove(key K) {
	currNode, parentNode := t.FinderNode(key)

	defer t.UnlockParent(parentNode)

	if currNode == nil {
		return
	}

	defer currNode.Unlock()

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

		defer tmpNode.Unlock()
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

func (oNd *OptimisticNode[T, K]) Lock() {
	if oNd == nil {
		return
	}
	oNd.mutex.Lock()
}

func (oNd *OptimisticNode[T, K]) Unlock() {
	if oNd == nil {
		return
	}
	oNd.mutex.Unlock()
}

func (t *OptimisticTree[T, K]) UnlockParent(parent *OptimisticNode[T, K]) {
	if parent == nil {
		t.mutex.Unlock()
	} else {
		parent.Unlock()
	}
}

func (t *OptimisticTree[T, K]) FinderNode(key K) (currentNode *OptimisticNode[T, K], parentNode *OptimisticNode[T, K]) {
	for {
		t.mutex.Lock()

		if t.root == nil {
			return
		}

		tmpNode := t.root
		var tmpPrevNode *OptimisticNode[T, K] = nil

		for tmpNode != nil && tmpNode.key != key {
			tmpGrandNode := tmpPrevNode
			tmpPrevNode = tmpNode

			switch cmp.Compare(key, tmpNode.key) {
			case -1:
				tmpNode = tmpNode.left
			case 1:
				tmpNode = tmpNode.right
			}
			if tmpGrandNode == nil {
				t.mutex.Unlock()
			}
		}

		if tmpPrevNode != nil {
			tmpPrevNode.Lock()
		}
		if tmpNode != nil {
			tmpNode.Lock()
		}

		if t.Validate(key, tmpNode, tmpPrevNode) {
			return tmpNode, tmpPrevNode
		}
		tmpNode.Unlock()
		tmpPrevNode.Unlock()
	}
}

func (t *OptimisticTree[T, K]) Validate(key K, curr, parent *OptimisticNode[T, K]) bool {
	if curr == nil && parent == nil {
		return t.root == nil
	}
	tmpNode := t.root
	var prevNode *OptimisticNode[T, K] = nil

	for tmpNode != nil && tmpNode.key != key && tmpNode != curr {
		prevNode = tmpNode
		switch cmp.Compare(key, tmpNode.key) {
		case -1:
			tmpNode = tmpNode.left
		case 1:
			tmpNode = tmpNode.right
		}
	}
	return curr == tmpNode && parent == prevNode
}

func (t *OptimisticTree[T, K]) IsValid() bool {
	return t.root.isValid()
}

func (oNd *OptimisticNode[T, K]) isValid() bool {
	if oNd == nil {
		return true
	}
	if oNd.left != nil && oNd.left.key >= oNd.key {
		return false
	}
	if oNd.right != nil && oNd.right.key <= oNd.key {
		return false
	}
	return oNd.left.isValid() && oNd.right.isValid()
}
