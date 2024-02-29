package benchmarks

import (
	"Treiber-stack/stacks"
	"sync"
	"testing"
)

const count_elem = 1_000_000

func Benchmark(b *testing.B) {
	simpleStack := stacks.CreateSimpleStack[int]()
	treiberStack := stacks.CreateTreiberStack[int]()

	b.Run("SimpleStack", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < count_elem; j++ {
				simpleStack.Push(j)
			}
			for j := 0; j < count_elem; j++ {
				_, err := simpleStack.Pop()
				if err != nil {
					b.Errorf("Unexpected error in %s stack: %d", "simple stack", err)
				}
			}
		}
	})

	b.Run("TreiberStack not concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < count_elem; j++ {
				treiberStack.Push(j)
			}
			for j := 0; j < count_elem; j++ {
				_, err := treiberStack.Pop()
				if err != nil {
					b.Errorf("Unexpected error in %s stack: %d", "treiber stack", err)
				}
			}
		}
	})

	b.Run("TreiberStack concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := sync.WaitGroup{}
			wg.Add(count_elem)
			for j := 0; j < count_elem; j++ {
				go func(j int) {
					treiberStack.Push(j)
					wg.Done()
				}(j)
			}
			wg.Wait()

			wg.Add(count_elem)
			for j := 0; j < count_elem; j++ {
				go func(j int) {
					_, err := treiberStack.Pop()
					if err != nil {
						b.Errorf("Unexpected error in %s stack: %d", "treiber stack", err)
					}
					wg.Done()
				}(j)
			}
			wg.Wait()
		}
	})
}
