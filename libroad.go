package frizzante

import (
	"strings"
)

type Road struct {
	lanes map[string]chan int
}

func RoadCreate() *Road {
	return &Road{lanes: map[string]chan int{}}
}

func RoadWithLane(self *Road, key ...string) chan int {
	path := strings.Join(key, ":")
	lane, laneExists := self.lanes[path]
	if laneExists {
		return lane
	}

	lane = make(chan int, 1)
	self.lanes[path] = lane
	self.lanes[path] <- 0
	return lane
}

func RoadWithoutLane(self *Road, key ...string) {
	path := strings.Join(key, ":")
	delete(self.lanes, path)
}
