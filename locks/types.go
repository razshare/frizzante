package locks

import "sync"

type Lock struct {
	Id string
	sync.Mutex
}
