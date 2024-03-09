package tests

import (
	"Treiber-stack/stacks"
	"Treiber-stack/stacks/Simple"
	"Treiber-stack/stacks/Treiber"
	"fmt"
	"sync"
	"testing"
)

func TestPopAndPush(t *testing.T) {
	simpleSt := Simple.CreateSimpleStack[int]()
	treiberSt := Treiber.CreateTreiberStack[int]()
	optTreiberSt := Treiber.CreateTreiberStack[int]()
	var tests = []struct {
		currStack stacks.Stack[int]
		typeStack string
	}{
		{&simpleSt, "simple"},
		{&treiberSt, "treiber"},
		{&optTreiberSt, "optimization treiber"},
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
	simpleSt := Simple.CreateSimpleStack[int]()
	treiberSt := Treiber.CreateTreiberStack[int]()
	optTreiberSt := Treiber.CreateTreiberStack[int]()
	var tests = []struct {
		currStack stacks.Stack[int]
		typeStack string
	}{
		{&simpleSt, "simple"},
		{&treiberSt, "treiber"},
		{&optTreiberSt, "optimization treiber"},
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
		if sz := myStack.Size(); sz != 100 {
			t.Errorf("Size expected %d, but get %d", 100, sz)
		}
	}
}

func TestPushGoroutines(t *testing.T) {
	treiberSt := Treiber.CreateTreiberStack[int]()
	optTreiberSt := Treiber.CreateTreiberStack[int]()
	var tests = []struct {
		currStack stacks.Stack[int]
		typeStack string
	}{
		{&treiberSt, "treiber"},
		{&optTreiberSt, "optimization treiber"},
	}
	for _, testStruct := range tests {
		myStack := testStruct.currStack
		goroutineCount := 100
		wg := sync.WaitGroup{}
		wg.Add(goroutineCount)
		for i := 0; i < goroutineCount; i++ {
			go func() {
				for j := 0; j < 10_000; j++ {
					myStack.Push(j)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		currSize := myStack.Size()
		fmt.Println("Size stack before all pop", currSize)
		if goroutineCount*10_000 != currSize {
			t.Errorf("Expected %d, but get %d", goroutineCount*10_000, currSize)
			return
		}

		wg.Add(goroutineCount)
		for i := 0; i < goroutineCount; i++ {
			go func() {
				for j := 0; j < 10_000; j++ {
					_, err := myStack.Pop()
					if err != nil {
						t.Errorf("%d", err)
						return
					}
				}
				wg.Done()
			}()
		}

		wg.Wait()
		currSize = myStack.Size()
		fmt.Println("Size stack after all pop", currSize)
	}
}
