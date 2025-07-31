package main

import (
	"github.com/razshare/frizzante/locks"
	"testing"
)

func TestNew(test *testing.T) {
	lock := locks.New()
	if len(lock.Names) > 0 {
		test.Fatal("road lanes should be empty")
	}
}

func TestLock(test *testing.T) {
	lock := locks.New()
	mutex := locks.FindAndAcquire(lock, "asd", "asd")
	if len(lock.Names) != 1 {
		test.Fatal("road should count exactly 1 lane")
	}
	if locks.FindAndAcquire(lock, "asd", "asd") != mutex {
		test.Fatal("mutexes should match")
	}
}

func TestRelease(test *testing.T) {
	lock := locks.New()
	mutex := locks.FindAndAcquire(lock, "asd", "asd")
	if len(lock.Names) != 1 {
		test.Fatal("road should count exactly 1 lane")
	}

	locks.Release(lock, "asd", "asd")

	if len(lock.Names) != 0 {
		test.Fatal("road should not contain any lanes")
	}

	if locks.FindAndAcquire(lock, "asd", "asd") == mutex {
		test.Fatal("mutexes should not match")
	}
}
