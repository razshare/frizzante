package main

import (
	"github.com/razshare/frizzante/locks"
	"testing"
)

func TestAcquire(test *testing.T) {
	locks.Acquire("test")

	// Cleanup.
	defer func() { delete(locks.Map, "test") }()

	if len(locks.Map) == 0 {
		test.Fatal("locks should count 0")
	}
}

func TestDestroy(test *testing.T) {
	lock := locks.Acquire("test")

	// Cleanup.
	defer func() { delete(locks.Map, "test") }()

	if len(locks.Map) != 1 {
		test.Fatal("parallels should count 1")
	}

	lock.Destroy()

	if len(locks.Map) != 0 {
		test.Fatal("parallels should count 0")
	}
}
