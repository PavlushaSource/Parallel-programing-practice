package benchmarks

import (
	"Treiber-stack/stacks"
	"Treiber-stack/stacks/Simple"
	"Treiber-stack/stacks/Treiber"
	"Treiber-stack/stacks/optimizationTreiber"
	"runtime"
	"sync"
	"testing"
)

const countElem = 1_000_000

func NonConcurrentPushAndPop(stack stacks.Stack[int]) {
	for j := 0; j < countElem; j++ {
		stack.Push(j)
	}
	for j := 0; j < countElem; j++ {
		stack.Pop()
	}
}

func BenchmarkNonConcurrent(b *testing.B) {
	b.Run("SimpleStack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			simpleStack := Simple.CreateSimpleStack[int]()
			NonConcurrentPushAndPop(&simpleStack)
		}
	})

	b.Run("TreiberStack not concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			NonConcurrentPushAndPop(&treiberStack)
		}
	})

	b.Run("Optimization back-off elimination treiberStack not concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			NonConcurrentPushAndPop(&optimizeTreiberStack)
		}
	})
}

func littleConcurrent(stack stacks.Stack[int]) {
	goroutineCount := 100
	wg := sync.WaitGroup{}
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func() {
			for j := 0; j < countElem/goroutineCount; j++ {
				stack.Push(j)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func() {
			for j := 0; j < countElem/goroutineCount; j++ {
				stack.Pop()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func allConcurrent(stack stacks.Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(countElem)
	for j := 0; j < countElem; j++ {
		go func(j int) {
			stack.Push(j)
			wg.Done()
		}(j)
	}
	wg.Wait()

	wg.Add(countElem)

	for j := 0; j < countElem; j++ {
		go func() {
			stack.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkConcurrent(b *testing.B) {
	runtime.GOMAXPROCS(12)

	b.Run("TreiberStack little concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			littleConcurrent(&treiberStack)
		}
	})

	b.Run("Optimization back-off elimination treiberStack little concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			littleConcurrent(&optimizeTreiberStack)
		}
	})

	b.Run("TreiberStack all concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			allConcurrent(&treiberStack)
		}
	})

	b.Run("TreiberStack all concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			allConcurrent(&optimizeTreiberStack)
		}
	})
}

func PushAndPopInRow(stack stacks.Stack[int]) {
	stack.Push(5)
	stack.Pop()
	stack.Push(20)
	stack.Pop()
	stack.Push(10)
	stack.Push(10)
	stack.Pop()
	stack.Pop()
	stack.Push(33)
}

func BenchmarkOptimizationCompare(b *testing.B) {
	runtime.GOMAXPROCS(12)

	b.Run("TreiberStack the standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			for j := 0; j < 100_000; j++ {
				go PushAndPopInRow(&treiberStack)
			}
		}
	})

	b.Run("TreiberStack with back-off elimination random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			for j := 0; j < 100_000; j++ {
				go func(j int) {
					for k := 0; k < 9; k++ {
						optimizeTreiberStack.Push(j)
					}
				}(j)
			}
		}
	})

	b.Run("TreiberStack with back-off elimination smart", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			for j := 0; j < 100_000; j++ {
				go PushAndPopInRow(&optimizeTreiberStack)
			}
		}
	})
}
