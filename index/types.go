package index

import (
	"embed"
	"github.com/razshare/frizzante/server"
	"log"
)

type Config struct {
	Development bool
	Document    string
	Root        string
	Channels    Channels
	Efs         embed.FS
	InfoLog     *log.Logger
	ErrorLog    *log.Logger
	Render      server.Render
}

type Channels struct {
	Document chan string
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
