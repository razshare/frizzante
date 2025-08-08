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
	"strings"
)

func ProduceScript(
	efs embed.FS,
	appRoot string,
	fileName string,
	infoLog *log.Logger,
	errorLog *log.Logger,
) *Script {
	var exit bool
	var stop = make(chan any, 1)
	var script = make(chan string, 1)
	var fileNameFixed = strings.ReplaceAll(fileName, "\\", "/")
	go func() { <-stop; exit = true }()
	go func() {
		defer infoLog.Println("script producer stopped")
		for !exit {
			data, readError := efs.ReadFile(fileNameFixed)
			if readError != nil {
				errorLog.Println(readError, stack.Trace())
				continue
			}
			bundledSourceCode, bundleError := js.Bundle(appRoot, api.FormatCommonJS, string(data))
			if bundleError != nil {
				errorLog.Println(bundleError, stack.Trace())
				continue
			}
			script <- fmt.Sprintf(globals.RenderScriptFormat, bundledSourceCode)
		}
	}()
	return &Script{
		Value: script,
		Stop:  stop,
	}
}

func ProduceDocument(
	efs embed.FS,
	fileName string,
	infoLog *log.Logger,
	errorLog *log.Logger,
) *Document {
	var exit bool
	var stop = make(chan any, 1)
	var document = make(chan string, 1)
	var fileNameFixed = strings.ReplaceAll(fileName, "\\", "/")
	go func() { <-stop; exit = true }()
	go func() {
		defer infoLog.Println("document producer stopped")
		for !exit {
			data, readError := efs.ReadFile(fileNameFixed)
			if readError != nil {
				errorLog.Println(readError, stack.Trace())
				continue
			}
			document <- string(data)
		}
	}()
	return &Document{
		Value: document,
		Stop:  stop,
	}
}

func ProduceRuntime(
	parallels uint32,
	infoLog *log.Logger,
) *Runtime {
	var exit bool
	var stop = make(chan any, 1)
	var runtime = make(chan *goja.Runtime, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer infoLog.Println("runtime producer stopped")
		var count uint32
		for !exit {
			if count >= parallels {
				infoLog.Printf("maximum number of parallel runtimes (%d) reached", parallels)
				return
			}
			runtime <- goja.New()
			count++
		}
	}()
	return &Runtime{
		Value: runtime,
		Stop:  stop,
	}
}

func ProduceProgram(
	scriptName string,
	script chan string,
	parallels uint32,
	infoLog *log.Logger,
	errorLog *log.Logger,
) *Program {
	var exit bool
	var stop = make(chan any, 1)
	var program = make(chan *goja.Program, 1)
	go func() { <-stop; exit = true }()
	go func() {
		defer infoLog.Println("program producer stopped")
		var count uint32
		for !exit {
			if count >= parallels {
				infoLog.Printf("maximum number of parallel programs (%d) reached", parallels)
				return
			}

			value, compileError := goja.Compile(scriptName, <-script, false)
			if compileError != nil {
				errorLog.Println(compileError, stack.Trace())
				continue
			}
			program <- value
			count++
		}
	}()
	return &Program{
		Value: program,
		Stop:  stop,
	}
}
