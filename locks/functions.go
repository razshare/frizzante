package locks

import (
	"strings"
	"sync"
)

func New() *Lock {
	return &Lock{Names: map[string]*sync.Mutex{}}
}

// FindAndAcquire adds a new lane to the road.
func FindAndAcquire(lock *Lock, keys ...string) *sync.Mutex {
	path := strings.Join(keys, ":")
	lane, laneExists := lock.Names[path]
	if laneExists {
		return lane
	}

	var newLane sync.Mutex
	lock.Names[path] = &newLane
	return &newLane
}

// Release unlocks the lane and removes it from the road.
func Release(lock *Lock, keys ...string) *Lock {
	path := strings.Join(keys, ":")
	delete(lock.Names, path)
	return lock
}
