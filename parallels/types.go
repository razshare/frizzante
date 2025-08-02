package parallels

type Complete[T any] = chan<- T
type Completed[T any] = <-chan T
type ComposeRunner[T any] = func(complete Complete[T])
type FirstRunner[T any] = func(complete Complete[T], completed Completed[T])

type Parallel[T any] struct {
	Value T
}
