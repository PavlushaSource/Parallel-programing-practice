package optimizationTreiber

import (
	"errors"
	"sync/atomic"
)

type exchangerState int

const (
	empty exchangerState = iota
	wait
	busy
)

type exchanger[T any] struct {
	item atomic.Value
}

type exchangeItem[T any] struct {
	value *T
	state exchangerState
}

func newExchanger[T any]() exchanger[T] {
	newEx := exchanger[T]{}
	newEx.item.Store(exchangeItem[T]{state: empty})
	return newEx
}

func (ex *exchanger[T]) exchange(val *T, waitSteps int) (*T, error) {
	emptyCase := func(passSteps int) (*T, error) {
		for j := passSteps; j < waitSteps; j++ {
			exItem := ex.item.Load().(exchangeItem[T])
			if exItem.state == busy {
				newItem := exchangeItem[T]{state: empty}
				ex.item.Store(newItem)
				return exItem.value, nil
			}
		}
		return new(T), errors.New("end cycle")
	}

	for i := 0; i < waitSteps; i++ {
		exItem := ex.item.Load().(exchangeItem[T])

		if exItem.state == empty {
			oldItem := exchangeItem[T]{state: empty}
			newItem := exchangeItem[T]{value: val, state: wait}
			if ex.item.CompareAndSwap(oldItem, newItem) {
				return emptyCase(i)
			}
		} else if exItem.state == wait {
			oldItem := exchangeItem[T]{value: exItem.value, state: wait}
			newItem := exchangeItem[T]{value: val, state: busy}
			if ex.item.CompareAndSwap(oldItem, newItem) {
				return exItem.value, nil
			}
		}
	}
	return new(T), errors.New("end cycle")
}
