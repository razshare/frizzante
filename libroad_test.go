package frizzante

import (
	"github.com/razshare/frizzante/libroad"
	"testing"
)

func TestRoad_NewRoad(t *testing.T) {
	road := libroad.NewRoad()
	if len(road.Lanes) > 0 {
		t.Fatal("road lanes should be empty")
	}
}

func TestRoad_WithLane(t *testing.T) {
	road := libroad.NewRoad()
	mut := road.WithLane("asd", "asd")
	if len(road.Lanes) != 1 {
		t.Fatal("road should count exactly 1 lane")
	}
	if road.WithLane("asd", "asd") != mut {
		t.Fatal("mutexes should match")
	}
}

func TestRoad_WithoutLane(t *testing.T) {

	road := libroad.NewRoad()
	mut := road.WithLane("asd", "asd")
	if len(road.Lanes) != 1 {
		t.Fatal("road should count exactly 1 lane")
	}

	road.WithoutLane("asd", "asd")

	if len(road.Lanes) != 0 {
		t.Fatal("road should not contain any lanes")
	}

	if road.WithLane("asd", "asd") == mut {
		t.Fatal("mutexes should not match")
	}
}
