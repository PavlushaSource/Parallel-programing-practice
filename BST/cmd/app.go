package main

import (
	"BST/trees"
	"fmt"
	"sync"
)

func main() {
	tree := trees.NewGrainedSyncTree[int, int]()
	wg := sync.WaitGroup{}
	wg.Wait()
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			tree.Insert(j, j*10)
		}(i)
	}
	wg.Wait()
	for i := 10; i > 0; i-- {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			v, exist := tree.Find(j)
			fmt.Printf("key = %d, value = %d, exist = %t\n", j, v, exist)
		}(i)
	}
	wg.Wait()
	fmt.Println("REMOVE ALL")
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			tree.Remove(j)
		}(i)
	}
	wg.Wait()

	for i := 10; i > 0; i-- {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			v, exist := tree.Find(j)
			fmt.Printf("key = %d, value = %d, exist = %t\n", j, v, exist)
		}(i)
	}
	wg.Wait()
}
