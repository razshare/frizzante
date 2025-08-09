package container

import (
	"embed"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/js"
	"github.com/razshare/frizzante/stack"
	"log"
	"time"
)

func ProduceScript(efs embed.FS, root string, name string, ilog *log.Logger, elog *log.Logger) *Script {
	var exit bool
	var stop = make(chan any, 1)
	var script = make(chan string, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer ilog.Println("script producer stopped")
		for !exit {
			d, err := efs.ReadFile(name)
			if err != nil {
				elog.Println(err, stack.Trace())
				time.Sleep(time.Second)
				continue
			}
			src, err := js.Bundle(root, api.FormatCommonJS, string(d))
			if err != nil {
				elog.Println(err, stack.Trace())
				time.Sleep(time.Second)
				continue
			}
			script <- fmt.Sprintf(globals.RenderScriptFormat, src)
		}
	}()
	return &Script{
		Value: script,
		Stop:  stop,
	}
}

func ProduceDocument(efs embed.FS, name string, ilog *log.Logger, elog *log.Logger) *Document {
	var exit bool
	var stop = make(chan any, 1)
	var doc = make(chan string, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer ilog.Println("document producer stopped")
		for !exit {
			d, err := efs.ReadFile(name)
			if err != nil {
				elog.Println(err, stack.Trace())
				time.Sleep(time.Second)
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

func ProduceRuntime(parallel uint32, ilog *log.Logger) *Runtime {
	var exit bool
	var stop = make(chan any, 1)
	var run = make(chan *goja.Runtime, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer ilog.Println("runtime producer stopped")
		var count uint32
		for !exit {
			if count >= parallel {
				ilog.Printf("maximum number of parallel runtimes (%d) reached", parallel)
				return
			}
			run <- goja.New()
			count++
		}
	}()
	return &Runtime{
		Value: run,
		Stop:  stop,
	}
}

func ProduceProgram(name string, script chan string, parallel uint32, ilog *log.Logger, elog *log.Logger) *Program {
	var exit bool
	var stop = make(chan any, 1)
	var prog = make(chan *goja.Program, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer ilog.Println("program producer stopped")
		var count uint32
		for !exit {
			if count >= parallel {
				ilog.Printf("maximum number of parallel programs (%d) reached", parallel)
				return
			}

			value, err := goja.Compile(name, <-script, false)
			if err != nil {
				elog.Println(err, stack.Trace())
				time.Sleep(time.Second)
				continue
			}
			prog <- value
			count++
		}
	}()
	return &Program{
		Value: prog,
		Stop:  stop,
	}
}
