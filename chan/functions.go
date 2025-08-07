package _chan

import "sync"

// Join joins multiple channels into one.
func Join[T any](ch ...chan T) chan T {
	channel := make(chan T, 1)

	go func() {
		var group sync.WaitGroup
		group.Add(len(ch))
		for _, channelN := range ch {
			go func() { defer group.Done(); channel <- <-channelN }()
		}
		group.Wait()
		close(channel)
	}()

	return channel
}
