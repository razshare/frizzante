package index

import (
	"embed"
	"github.com/razshare/frizzante/view"
	"log"
)

type Index struct {
	Development bool
	Document    string
	Root        string
	Channels    Channels
	Efs         embed.FS
	ErrorLog    *log.Logger
	Render      func(v view.View) (string, error)
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
