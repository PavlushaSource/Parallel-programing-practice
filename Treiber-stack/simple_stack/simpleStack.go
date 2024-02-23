package simple_stack

import (
	"errors"
)

type SimpleStack[T any] struct {
	head *Node[T]
}

type Node[T any] struct {
	value T
	next  *Node[T]
}

func (stack *SimpleStack[T]) Peek() T {
	return stack.head.value
}

func (stack *SimpleStack[T]) Pop() (T, error) {
	if stack.head == nil {
		var nilVal T
		return nilVal, errors.New("stack is already empty")
	}
	lastValue := stack.head.value
	stack.head = stack.head.next
	return lastValue, nil
}

func (stack *SimpleStack[T]) Push(val T) {
	newNode := Node[T]{value: val}
	stack.head, newNode.next = &newNode, stack.head
}

func CreateSimpleStack[T any]() SimpleStack[T] {
	return SimpleStack[T]{}
}

type Stack[T any] interface {
	Push(T)
	Pop() (T, error)
	Peek() T
}

//func main() {
//	st := CreateSimpleStack[int]()
//	st.Push(10)
//	st.Push(20)
//	fmt.Println(st.Peek())
//	val, err := st.Pop()
//	if err != nil {
//		fmt.Println(val, err)
//	} else {
//		fmt.Println(val)
//	}
//	fmt.Println(st.Peek())
//}
