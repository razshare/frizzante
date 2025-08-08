package container

import (
	"embed"
	"github.com/dop251/goja"
	"log"
	"sync"
)

type Config struct {
	PublicRoot  string
	Development bool
	Parallels   uint32
	Document    string
	Script      string
	Root        string
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	Efs         embed.FS
}

type Channels struct {
	Document chan string
	Script   chan string
	Program  chan *goja.Program
	Runtime  chan *goja.Runtime
	Stop     chan any
}

type Container struct {
	Channels Channels
	Config   Config
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
