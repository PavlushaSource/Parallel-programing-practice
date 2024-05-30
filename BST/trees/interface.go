package trees

import (
	"cmp"
)

type Tree[T any, K cmp.Ordered] interface {
	Find(K) (T, bool)
	Insert(K, T)
	Remove(K)
	IsValid() bool
}
