package frizzante

type Writable[T any] struct {
	subscribers map[int]func(value T)
	index       int
	value       T
}

// WritableCreate creates a writable store.
func WritableCreate[T any](value T) *Writable[T] {
	return &Writable[T]{
		value:       value,
		subscribers: map[int]func(value T){},
	}
}

// WritableSubscribe subscribes to a writable store
// and returns a function that removes the subscription.
//
// You should always remove the subscription when you're done with the store.
func WritableSubscribe[T any](self *Writable[T], callback func(value T)) (unsubscribe func()) {
	index := self.index
	self.index++
	self.subscribers[index] = callback
	unsubscribe = func() {
		delete(self.subscribers, index)
	}
	return
}

// WritableRead reads the current value of a writable store.
func WritableRead[T any](self *Writable[T]) T {
	return self.value
}

// WritableWrite writes a new value to the writable store.
func WritableWrite[T any](self *Writable[T], value T) {
	self.value = value
	for _, subscriber := range self.subscribers {
		subscriber(value)
	}
}
