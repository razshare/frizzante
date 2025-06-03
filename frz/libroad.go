package frz

import (
	"strings"
	"sync"
)

type Road struct {
	lanes map[string]*sync.Mutex
}

func NewRoad() *Road {
	return &Road{lanes: map[string]*sync.Mutex{}}
}

// WithLane adds a new lane to the road.
func (road *Road) WithLane(key ...string) *sync.Mutex {
	path := strings.Join(key, ":")
	lane, laneExists := road.lanes[path]
	if laneExists {
		return lane
	}

	var newLane sync.Mutex
	road.lanes[path] = &newLane
	return &newLane
}

// WithoutLane unlocks the lane and removes it from the road.
func (road *Road) WithoutLane(key ...string) *Road {
	path := strings.Join(key, ":")
	delete(road.lanes, path)
	return road
}
