package main

import (
	"github.com/razshare/frizzante/locks"
	"testing"
)

func TestAcquire(test *testing.T) {
	locks.Acquire("test")

	// Cleanup.
	defer func() { delete(locks.Map, "test") }()

	if locks.Map["test"] == nil {
		test.Fatal("locks should not contain test")
	}
}

func TestDestroy(test *testing.T) {
	lock := locks.Acquire("test")

	// Cleanup.
	defer func() { delete(locks.Map, "test") }()

	if locks.Map["test"] == nil {
		test.Fatal("locks should not contain test")
	}

	lock.Destroy()

	if locks.Map["test"] != nil {
		test.Fatal("locks should contain test")
	}
}
