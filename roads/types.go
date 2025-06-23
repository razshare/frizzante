package roads

import "sync"

type Road struct {
	Lanes map[string]*sync.Mutex
}
