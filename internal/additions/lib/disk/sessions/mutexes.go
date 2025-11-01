package sessions

import "sync"

var Mutexes = map[string]*sync.Mutex{}
