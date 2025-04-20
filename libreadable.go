package frizzante

type Readable[T any] struct {
	subscribers map[int]func(value T)
	index       int
	destroy     func()
	value       T
}

// ReadableCreate creates a readable store.
func ReadableCreate[T any](
	value T,
	update func(
		set func(value T),
	) func(),
) *Readable[T] {
	readable := &Readable[T]{
		value:       value,
		subscribers: map[int]func(value T){},
	}

	readable.destroy = update(func(value T) {
		readable.value = value
		for _, subscriber := range readable.subscribers {
			subscriber(value)
		}
	})

	return readable
}

// ReadableSubscribe subscribes to a readable store
// and returns a function that removes the subscription.
//
// You should always remove the subscription when you're done with the store.
func ReadableSubscribe[T any](self *Readable[T], callback func(value T)) (unsubscribe func()) {
	index := self.index
	self.index++
	self.subscribers[index] = callback
	unsubscribe = func() {
		delete(self.subscribers, index)
		if 0 == len(self.subscribers) {
			self.destroy()
		}
	}
	return
}

// ReadableGet get the current value of a readable store.
func ReadableGet[T any](self *Readable[T]) T {
	return self.value
}
