package index

import (
	"embed"
	"github.com/razshare/frizzante/stack"
	"log"
	"time"
)

func ProduceDocument(efs embed.FS, name string, elog *log.Logger) *Document {
	var exit bool
	var stop = make(chan any, 1)
	var doc = make(chan string, 1)
	go func() { <-stop; exit = true }()
	go func() {
		for !exit {
			d, err := efs.ReadFile(name)
			if err != nil {
				elog.Println(err, stack.Trace())
				time.Sleep(10 * time.Second)
				continue
			}
			doc <- string(d)
		}
	}()
	return &Document{
		Value: doc,
		Stop:  stop,
	}
}
