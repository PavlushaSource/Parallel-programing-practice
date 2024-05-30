package benchmarks

import (
	"BST/trees"
	"sync"
	"testing"
)

const countElem = 10_000

func SeqInsert(t trees.Tree[int, int]) {
	for i := 0; i < countElem; i++ {
		t.Insert(i, i)
	}
}

func SeqRemove(t trees.Tree[int, int]) {
	for i := 0; i < countElem; i++ {
		t.Remove(i)
	}
}

func BenchmarkSeqInsert(b *testing.B) {
	b.Run("Grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewGrainedSyncTree[int, int]()
			SeqInsert(tree)
		}
	})

	b.Run("Fine-grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewFineGrainedSyncTree[int, int]()
			SeqInsert(tree)
		}
	})

	b.Run("Optimistic Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewOptimisticSyncTree[int, int]()
			SeqInsert(tree)
		}
	})
}

func BenchmarkSeqRemove(b *testing.B) {
	b.Run("Grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewGrainedSyncTree[int, int]()
			SeqRemove(tree)
		}
	})

	b.Run("Fine-grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewFineGrainedSyncTree[int, int]()
			SeqRemove(tree)
		}
	})

	b.Run("Optimistic Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewOptimisticSyncTree[int, int]()
			SeqRemove(tree)
		}
	})
}

func ConcurrentInsert(t trees.Tree[int, int], wg *sync.WaitGroup) {
	for i := 0; i < countElem*10; i++ {
		go func(key, value int) {
			defer wg.Done()
			t.Insert(key, value)
		}(i, i)
	}
}

func ConcurrentRemove(t trees.Tree[int, int], wg *sync.WaitGroup) {
	for i := 0; i < countElem*10; i++ {
		go func(key int) {
			defer wg.Done()
			t.Remove(key)
		}(i)
	}
}

func BenchmarkConcurrentInsertAndRemove(b *testing.B) {
	b.Run("Grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewGrainedSyncTree[int, int]()
			wg := sync.WaitGroup{}
			wg.Add(countElem * 10 * 2)
			go ConcurrentInsert(tree, &wg)
			go ConcurrentRemove(tree, &wg)
			wg.Wait()
		}
	})

	b.Run("Fine-grained Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewFineGrainedSyncTree[int, int]()
			wg := sync.WaitGroup{}
			wg.Add(countElem * 10 * 2)
			go ConcurrentInsert(tree, &wg)
			go ConcurrentRemove(tree, &wg)
			wg.Wait()
		}
	})

	b.Run("Optimistic Tree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tree := trees.NewOptimisticSyncTree[int, int]()
			wg := sync.WaitGroup{}
			wg.Add(countElem * 10 * 2)
			go ConcurrentInsert(tree, &wg)
			go ConcurrentRemove(tree, &wg)
			wg.Wait()
		}
	})
}
