package simple_stack

import "testing"

func TestPopAndPush(t *testing.T) {
	simpleSt := CreateSimpleStack[int]()
	var tests = []struct {
		currStack Stack[int]
		typeStack string
	}{
		{&simpleSt, "simple"},
	}

	for _, testStruct := range tests {
		myStack := testStruct.currStack
		elements := 100
		for i := 0; i < elements; i++ {
			myStack.Push(i)
		}

		for i := 0; i < elements; i++ {
			res, err := myStack.Pop()
			if err != nil {
				t.Errorf("Unexpected error in %s stack: %d", testStruct.typeStack, err)
			} else if res != elements-1-i {
				t.Errorf("Expected %d on top of %s stack, but get %d", elements-1-i, testStruct.typeStack, res)
			}
		}

		_, err := myStack.Pop()
		if err == nil {
			t.Error("Stack expected to be empty")
		}
	}
}

func TestPush(t *testing.T) {
	simpleSt := CreateSimpleStack[int]()
	var tests = []struct {
		currStack Stack[int]
		typeStack string
	}{
		{&simpleSt, "simple"},
	}

	for _, testStruct := range tests {
		myStack := testStruct.currStack
		elements := 100
		for i := 0; i < elements; i++ {
			myStack.Push(i)
			res := myStack.Peek()
			if res != i {
				t.Errorf("Expected %d on top of %s stack, but get %d", elements-1-i, testStruct.typeStack, res)
			}
		}
	}
}
