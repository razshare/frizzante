package _chan

import "sync"

// Join joins multiple channels into one.
func Join[T any](channels ...chan T) chan T {
	channel := make(chan T, 1)

	go func() {
		var group sync.WaitGroup
		group.Add(len(channels))
		for _, channelN := range channels {
			go func() { defer group.Done(); channel <- <-channelN }()
		}
		group.Wait()
		close(channel)
	}()

	return channel
}
