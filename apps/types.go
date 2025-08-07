package apps

import (
	"github.com/dop251/goja"
	"log"
	"sync"
)

type Configuration struct {
	Development bool
	Parallels   uint32
	Document    string
	Script      string
	Root        string
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
}

type App struct {
	Document chan string
	Script   chan string
	Program  chan *goja.Program
	Runtime  chan *goja.Runtime
	Stop     chan any
}

type Script struct {
	Value chan string
	Stop  chan any
}

type Document struct {
	Value chan string
	Stop  chan any
}

type Program struct {
	Value chan *goja.Program
	Mutex sync.Mutex
	Stop  chan any
}

type Runtime struct {
	Value chan *goja.Runtime
	Mutex sync.Mutex
	Stop  chan any
}
