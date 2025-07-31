package locks

import (
	"strings"
	"sync"
)

func New() *Lock {
	return &Lock{Names: map[string]*sync.Mutex{}}
}

// Acquire finds a mutex and acquires it.
func (lock *Lock) Acquire(keys ...string) *sync.Mutex {
	path := strings.Join(keys, ":")
	lane, laneExists := lock.Names[path]
	if laneExists {
		return lane
	}

	var newLane sync.Mutex
	lock.Names[path] = &newLane
	return &newLane
}

// Remove releases.
func (lock *Lock) Remove(keys ...string) *Lock {
	path := strings.Join(keys, ":")
	delete(lock.Names, path)
	return lock
}
