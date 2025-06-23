package libroad

import (
	"strings"
	"sync"
)

func NewRoad() *Road {
	return &Road{Lanes: map[string]*sync.Mutex{}}
}

// WithLane adds a new lane to the road.
func (road *Road) WithLane(key ...string) *sync.Mutex {
	path := strings.Join(key, ":")
	lane, laneExists := road.Lanes[path]
	if laneExists {
		return lane
	}

	var newLane sync.Mutex
	road.Lanes[path] = &newLane
	return &newLane
}

// WithoutLane unlocks the lane and removes it from the road.
func (road *Road) WithoutLane(key ...string) *Road {
	path := strings.Join(key, ":")
	delete(road.Lanes, path)
	return road
}
