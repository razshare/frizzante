package locks

import (
	"strings"
	"sync"
)

func New() *Lock {
	return &Lock{Names: map[string]*sync.Mutex{}}
}

// Lock adds a new lane to the road.
func (road *Lock) Lock(keys ...string) *sync.Mutex {
	path := strings.Join(keys, ":")
	lane, laneExists := road.Names[path]
	if laneExists {
		return lane
	}

	var newLane sync.Mutex
	road.Names[path] = &newLane
	return &newLane
}

// Release unlocks the lane and removes it from the road.
func (road *Lock) Release(keys ...string) *Lock {
	path := strings.Join(keys, ":")
	delete(road.Names, path)
	return road
}
