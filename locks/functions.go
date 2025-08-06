package locks

import (
	"strings"
	"sync"
)

var Map = map[string]*Lock{}
var Mutex = &sync.Mutex{}

// Acquire retrieves or creates a Lock with a given key.
//
// Lock embeds Mutex.
//
// If the Lock is created anew, locks.Acquire locks the Lock immediately.
//
// Acquire itself guarantees thread-safety by locking to locks.Mutex.
//
// New locks are saved in locks.Map.
// Read Lock.Destroy for more details on the usage of locks.Map.
func Acquire(key ...string) *Lock {
	if len(key) == 0 {
		return &Lock{Mutex: sync.Mutex{}}
	}

	id := strings.Join(key, ":")

	Mutex.Lock()
	defer func() { Mutex.Unlock() }()

	lock, exists := Map[id]
	if exists {
		return lock
	}

	lock = &Lock{Id: id, Mutex: sync.Mutex{}}

	Map[id] = lock

	return lock
}

// Destroy removes the Lock from locks.Map.
//
// In order to guarantee thread-safety, Lock.Destroy
// locks both the Lock and locks.Mutex during this process.
//
// Both Lock and locks.Mutex are unlocked after
// the Lock has been removed from locks.Map.
func (lock *Lock) Destroy() {
	Mutex.Lock()
	lock.Mutex.Lock()

	defer func() {
		Mutex.Unlock()
		lock.Mutex.Unlock()
	}()

	delete(Map, lock.Id)
}
