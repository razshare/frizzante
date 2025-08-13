package app

import (
	"github.com/dop251/goja"
	"github.com/razshare/frizzante/server"
	"sync"
)

type Config struct {
	Development bool
	Parallels   uint32
	Document    string
	Script      string
	Root        string
	Channels    Channels
	Server      *server.Config
}

type Channels struct {
	Document chan string
	Script   chan string
	Program  chan *goja.Program
	Runtime  chan *goja.Runtime
	Stop     chan any
}

type App struct {
	Config Config
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
