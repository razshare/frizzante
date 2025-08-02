package locks

import (
	"strings"
	"sync"
)

var Map = map[string]*Lock{}
var Mutex = &sync.Mutex{}

// Acquire acquires a LockSynchronizer from ParallelMap
// based on the given key and returns its Mutex.
//
// If the LockSynchronizer doesn't exist, Lock creates it
// and if the key is not empty it also saves it in Map.
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

// Destroy removes the LockSynchronizer from ParallelMap.
func (lock *Lock) Destroy() {
	Mutex.Lock()
	lock.Mutex.Lock()

	defer func() {
		Mutex.Unlock()
		lock.Mutex.Unlock()
	}()

	delete(Map, lock.Id)
}
