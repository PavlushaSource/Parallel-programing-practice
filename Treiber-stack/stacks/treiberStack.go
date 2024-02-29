package stacks

import (
	"errors"
	"sync/atomic"
)

type TreiberStack[T any] struct {
	head atomic.Pointer[TNode[T]]
}

type TNode[T any] struct {
	value T
	next  atomic.Pointer[TNode[T]]
}

func (stack *TreiberStack[T]) Pop() (nilVar T, Err error) {
	for {
		head := stack.head.Load()
		if head == nil {
			return nilVar, errors.New("nil pointer to stack")
		}
		if stack.head.CompareAndSwap(head, head.next.Load()) {
			return head.value, nil
		}
	}
}

func (stack *TreiberStack[T]) Push(val T) {
	newHead := TNode[T]{value: val}
	for {
		head := stack.head.Load()
		newHead.next = atomic.Pointer[TNode[T]]{}
		newHead.next.Store(head)
		if stack.head.CompareAndSwap(head, &newHead) {
			return
		}
	}
}

func (stack *TreiberStack[T]) Peek() (nilVar T) {
	if stack == nil {
		return
	}
	head := stack.head.Load()
	if head == nil {
		return
	}
	return head.value
}

func CreateTreiberStack[T any]() TreiberStack[T] {
	return TreiberStack[T]{}
}

func (stack *TreiberStack[T]) Size() int {
	if stack == nil || stack.head.Load() == nil {
		return 0
	}
	elemCounter := 0
	curHead := stack.head.Load()
	for curHead != nil {
		elemCounter++
		curHead = curHead.next.Load()
	}
	return elemCounter
}
