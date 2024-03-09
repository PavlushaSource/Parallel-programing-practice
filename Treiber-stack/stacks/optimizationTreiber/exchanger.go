package optimizationTreiber

import "sync/atomic"

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

func (ex *exchanger[T]) exchange(val *T, waitSteps int) *T {
	for i := 0; i < waitSteps; i++ {
		item := ex.item.Load().(exchangeItem[T])

		if item.state == empty {
			newItem := exchangeItem[T]{value: val, state: wait}
			if ex.item.CompareAndSwap(item, newItem) {
				for i < waitSteps {
					item := ex.item.Load().(exchangeItem[T])
					if item.state == busy {
						newItem := exchangeItem[T]{state: empty}
						ex.item.Store(newItem)
						return item.value
					}
				}
			}
		} else if item.state == wait {
			newItem := exchangeItem[T]{value: val, state: busy}
			if ex.item.CompareAndSwap(item, newItem) {
				return item.value
			}
		}
	}
	return new(T)
}
