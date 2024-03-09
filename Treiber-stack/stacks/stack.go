package stacks

type Stack[T any] interface {
	Push(T)
	Pop() (T, error)
	Peek() T
	Size() int
}
