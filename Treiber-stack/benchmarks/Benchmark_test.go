package benchmarks

import (
	"Treiber-stack/stacks"
	"Treiber-stack/stacks/Simple"
	"Treiber-stack/stacks/Treiber"
	"Treiber-stack/stacks/optimizationTreiber"
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
			for j := 0; j < (countElem / goroutineCount); j++ {
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

func BenchmarkLittleConcurrent(b *testing.B) {
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
}

func BenchmarkAllConcurrent(b *testing.B) {
	b.Run("TreiberStack all concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			allConcurrent(&treiberStack)
		}
	})

	b.Run("Optimization back-off elimination treiberStack all concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			allConcurrent(&optimizeTreiberStack)
		}
	})
}

func PushAndPopInRow(stack stacks.Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(1_000)
	for j := 0; j < 1_000; j++ {
		go func(j int) {
			defer wg.Done()
			for k := 0; k < 1_000; k++ {
				stack.Push(j)
				stack.Pop()
			}
		}(j)
	}
	wg.Wait()
}

func PushPopConcurentRand(s stacks.Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(2_000)
	for j := 0; j < 1_000; j++ {
		go func(j int) {
			defer wg.Done()
			for k := 0; k < 1_000; k++ {
				s.Push(j)
			}
		}(j)
		go func() {
			wg.Done()
			for k := 0; k < 1_000; k++ {
				_, err := s.Pop()
				if err != nil {
					continue
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkOptimizationCompare(b *testing.B) {

	b.Run("TreiberStack push and pop in row", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			PushAndPopInRow(&treiberStack)
		}
	})

	b.Run("TreiberStack with back-off elimination push and pop in row", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			PushAndPopInRow(&optimizeTreiberStack)
		}
	})

	b.Run("TreiberStack random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := Treiber.CreateTreiberStack[int]()
			PushPopConcurentRand(&treiberStack)
		}
	})

	b.Run("TreiberStack with back-off elimination random", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := optimizationTreiber.CreateBackoffTreiberStack[int]()
			PushPopConcurentRand(&optimizeTreiberStack)
		}
	})
}
