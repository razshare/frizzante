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
	mutex := lock.Lock("asd", "asd")
	if len(lock.Names) != 1 {
		test.Fatal("road should count exactly 1 lane")
	}
	if lock.Lock("asd", "asd") != mutex {
		test.Fatal("mutexes should match")
	}
}

func TestRelease(test *testing.T) {
	lock := locks.New()
	mutex := lock.Lock("asd", "asd")
	if len(lock.Names) != 1 {
		test.Fatal("road should count exactly 1 lane")
	}

	lock.Release("asd", "asd")

	if len(lock.Names) != 0 {
		test.Fatal("road should not contain any lanes")
	}

	if lock.Lock("asd", "asd") == mutex {
		test.Fatal("mutexes should not match")
	}
}
