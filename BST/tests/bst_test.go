package tests

import (
	"BST/trees"
	"testing"
)

func TestInsert(t *testing.T) {
	grTree := trees.NewGrainedSyncTree[int, int]()
	fnGrTree := trees.NewFineGrainedSyncTree[int, int]()
	optTree := trees.NewOptimisticSyncTree[int, int]()

	var tests = []struct {
		currTree trees.Tree[int, int]
		typeSync string
	}{
		{grTree, "simple"},
		{fnGrTree, "fine grained"},
		{optTree, "optimistic"},
	}

	for _, testStruct := range tests {
		myTree := testStruct.currTree
		for i := 0; i < 100; i++ {
			myTree.Insert(i, i)
			if myTree.IsValid() != true {
				t.Errorf("Tree %s is not valid", testStruct.typeSync)
			}
		}
		for i := 99; i >= 0; i-- {
			if value, exist := myTree.Find(i); value != i || !exist {
				t.Errorf("Expected %d on top of %s tree, but get %d", i, testStruct.typeSync, value)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	grTree := trees.NewGrainedSyncTree[int, int]()
	fnGrTree := trees.NewFineGrainedSyncTree[int, int]()
	optTree := trees.NewOptimisticSyncTree[int, int]()

	var tests = []struct {
		currTree trees.Tree[int, int]
		typeSync string
	}{
		{grTree, "simple"},
		{fnGrTree, "fine grained"},
		{optTree, "optimistic"},
	}

	for _, testStruct := range tests {
		myTree := testStruct.currTree
		for i := 0; i < 100; i++ {
			myTree.Insert(i, i)
		}

		for i := 0; i < 100; i += 10 {
			if value, exist := myTree.Find(i); value != i || !exist {
				t.Errorf("Expected %d on top of %s tree, but get %d", i, testStruct.typeSync, value)
			}
		}

		for i := 0; i < 1000; i++ {
			myTree.Remove(i)
			if myTree.IsValid() != true {
				t.Errorf("Tree %s is not valid", testStruct.typeSync)
			}
		}

		for i := 0; i < 100; i++ {
			if _, exist := myTree.Find(i); exist {
				t.Errorf("Not remove element %d from %s tree", i, testStruct.typeSync)
			}
		}
	}
}
