package parallels

import "sync"

// Compose runs multiple functions in parallel and wait for each one to complete.
func Compose[T any](runners ...ComposeRunner[T]) []T {
	count := len(runners)
	if count == 0 {
		return make([]T, 0)
	}

	complete := make(chan T)
	values := make([]T, count)

	var secondary sync.WaitGroup
	secondary.Add(1)

	go func() {
		defer secondary.Done()
		index := 0
		for {
			value, more := <-complete
			if !more {
				return
			}

			if index < count {
				values[index] = value
				index++
			} else {
				values = append(values, value)
			}

		}
	}()

	var main sync.WaitGroup
	main.Add(count)

	for _, callback := range runners {
		go func() {
			defer main.Done()
			callback(complete)
		}()
	}

	main.Wait()

	close(complete)

	secondary.Wait()

	return values
}

// First runs multiple functions in parallel and wait for the first one to complete.
func First[T any](runners ...FirstRunner[T]) T {
	if len(runners) == 0 {
		var value T
		return value
	}

	complete := make(chan T)
	completed := make(chan T)

	for _, callback := range runners {
		go func() {
			callback(complete, completed)
		}()
	}

	value := <-complete
	completed <- value
	close(complete)

	return value
}
