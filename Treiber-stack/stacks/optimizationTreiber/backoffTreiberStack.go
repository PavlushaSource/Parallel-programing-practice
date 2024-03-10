package optimizationTreiber

import (
	"errors"
	"sync/atomic"
)

type OptimizedTreiberStack[T any] struct {
	head             atomic.Pointer[OTNode[T]]
	eliminationArray eliminationArray[T]
}

type OTNode[T any] struct {
	value T
	next  atomic.Pointer[OTNode[T]]
}

func (stack *OptimizedTreiberStack[T]) TryPop() (nilVar *T, Err error) {
	head := stack.head.Load()
	if head == nil {
		return nilVar, errors.New("nil pointer to stack")
	}
	if stack.head.CompareAndSwap(head, head.next.Load()) {
		return &head.value, nil
	}
	return
}

func (stack *OptimizedTreiberStack[T]) Pop() (nilVar T, Err error) {
	for {
		val, err := stack.TryPop()
		if err != nil {
			return nilVar, err
		}
		if val != nil {
			return *val, nil
		}
		valVisit, err := stack.eliminationArray.visit(nil)
		if val != nil && valVisit == nil {
			return *valVisit, nil
		}

	}
}

func (stack *OptimizedTreiberStack[T]) tryPush(n *OTNode[T]) bool {
	head := stack.head.Load()
	n.next.Store(head)
	return stack.head.CompareAndSwap(head, n)
}

func (stack *OptimizedTreiberStack[T]) Push(val T) {
	newHead := OTNode[T]{value: val}
	for {
		if stack.tryPush(&newHead) {
			return
		}
		valVisit, err := stack.eliminationArray.visit(&val)
		if valVisit == nil && err == nil {
			return
		}
	}
}

func (stack *OptimizedTreiberStack[T]) Peek() (nilVar T) {
	if stack == nil {
		return
	}
	head := stack.head.Load()
	if head == nil {
		return
	}
	return head.value
}

func (stack *OptimizedTreiberStack[T]) Size() int {
	elemCounter := 0
	if stack == nil || stack.head.Load() == nil {
		return 0
	}
	currHead := stack.head.Load()
	for currHead != nil {
		elemCounter++
		currHead = currHead.next.Load()
	}
	return elemCounter
}

func CreateBackoffTreiberStack[T any]() OptimizedTreiberStack[T] {
	return OptimizedTreiberStack[T]{eliminationArray: newEliminationArray[T](10, 1000)}
}
