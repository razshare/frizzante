package container

import (
	"embed"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/js"
	"github.com/razshare/frizzante/stack"
	"log"
	"os"
	"strings"
	"sync"
)

// Start starts a new app.
func Start(config Configuration, efs embed.FS) *Container {
	script := ProduceScript(
		efs,
		config.Root,
		config.Script,
		config.InfoLog,
		config.ErrorLog,
	)

	document := ProduceDocument(
		efs,
		config.Document,
		config.InfoLog,
		config.ErrorLog,
	)

	runtime := ProduceRuntime(
		config.Parallels,
		config.InfoLog,
	)

	program := ProduceProgram(
		config.Script,
		script.Value,
		config.Parallels,
		config.InfoLog,
		config.ErrorLog,
	)

	stop := make(chan any, 1)

	go func() {
		<-stop

		var group sync.WaitGroup
		group.Add(4)
		go func() { group.Done(); document.Stop <- 0 }()
		go func() { group.Done(); script.Stop <- 0 }()
		go func() { group.Done(); program.Stop <- 0 }()
		go func() { group.Done(); runtime.Stop <- 0 }()
		group.Wait()
	}()

	return &Container{
		Script:   script.Value,
		Document: document.Value,
		Program:  program.Value,
		Runtime:  runtime.Value,
		Stop:     stop,
	}
}

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

// ExecuteServerJs executes the server script.
func ExecuteServerJs(application *Container, configuration Configuration, properties map[string]any) (string, string, error) {
	var runtime *goja.Runtime
	var program *goja.Program
	var compileError error

	if configuration.Development {
		runtime = goja.New()
		var fileNameFixed = strings.ReplaceAll(configuration.Script, "\\", "/")
		data, rer := os.ReadFile(fileNameFixed)
		if rer != nil {
			return "", "", rer
		}
		source, bundleError := js.Bundle(configuration.Root, api.FormatCommonJS, string(data))
		if bundleError != nil {
			return "", "", bundleError
		}
		program, compileError = goja.Compile(configuration.Script, fmt.Sprintf(globals.RenderScriptFormat, source), false)
		if compileError != nil {
			return "", "", compileError
		}
	} else {
		runtime = <-application.Runtime
		program = <-application.Program
		defer func() { go func() { application.Runtime <- runtime }() }()
		defer func() { go func() { application.Program <- program }() }()
	}

	programResult, programError := runtime.RunProgram(program)
	if programError != nil {
		return "", "", programError
	}

	render, isFunction := goja.AssertFunction(programResult)

	if !isFunction {
		return "", "", errors.New("render is not a function")
	}

	promise, programError := render(goja.Undefined(), runtime.ToValue(properties))

	if programError != nil {
		return "", "", programError
	}

	value := promise.Export().(*goja.Promise).Result().ToObject(runtime)

	headValue := value.Get("head")
	bodyValue := value.Get("body")

	var head string
	var body string

	if headValue != nil {
		head = headValue.String()
	}

	if bodyValue != nil {
		body = bodyValue.String()
	}

	return head, body, nil
}
