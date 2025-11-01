package sessions

import (
	"path/filepath"
	"sync"
)

var Sessions = map[string]*Session{}
var Mutexes = map[string]*sync.Mutex{}
var DirectoryName = filepath.Join(".gen", "sessions")
