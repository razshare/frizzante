package receive

import "sync"

var Mutexes = map[string]*sync.Mutex{}
