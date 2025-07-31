package locks

import "sync"

type Lock struct {
	Names map[string]*sync.Mutex
}
