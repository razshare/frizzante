package frizzante

import (
	"strings"
)

type Road struct {
	lanes map[string]chan int
}

func NewRoad() *Road {
	return &Road{lanes: map[string]chan int{}}
}

func (road *Road) WithLane(key ...string) chan int {
	path := strings.Join(key, ":")
	lane, laneExists := road.lanes[path]
	if laneExists {
		return lane
	}

	lane = make(chan int, 1)
	road.lanes[path] = lane
	road.lanes[path] <- 0
	return lane
}

func (road *Road) WithoutLane(key ...string) *Road {
	path := strings.Join(key, ":")
	delete(road.lanes, path)
	return road
}
