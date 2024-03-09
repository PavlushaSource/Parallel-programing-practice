package optimizationTreiber

import "math/rand"

type eliminationArray[T any] struct {
	cap, waitSteps int
	exchangers     []exchanger[T]
}

func (elArr *eliminationArray[T]) visit(value *T) *T {
	index := rand.Intn(elArr.cap)
	return elArr.exchangers[index].exchange(value, elArr.waitSteps)
}

func newEliminationArray[T any](cap, waitSteps int) eliminationArray[T] {
	newArr := eliminationArray[T]{cap: cap, waitSteps: waitSteps}
	newArr.exchangers = make([]exchanger[T], cap)
	for i := range newArr.exchangers {
		newArr.exchangers[i] = newExchanger[T]()
	}
	return newArr
}
