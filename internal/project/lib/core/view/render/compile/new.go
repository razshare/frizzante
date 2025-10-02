package compile_

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func New(config Config) (render goja.Callable, runtime *goja.Runtime, err error) {
	var builder strings.Builder
	var server = filepath.Join(config.App, "dist", "app.server.js")
	var index = filepath.Join(config.App, "dist", "client", "index.html")

	server = strings.ReplaceAll(server, "\\", "/")
	index = strings.ReplaceAll(index, "\\", "/")

	runtime = goja.New()
	console := runtime.NewObject()
	createLogger := func(level LogLevel) func(call goja.FunctionCall) goja.Value {
		var logger *log.Logger

		switch level {
		case LogLevelDanger:
			logger = config.ErrorLog
		default:
			logger = config.InfoLog
		}

		return func(call goja.FunctionCall) goja.Value {
			builder.Reset()
			i := 0
			for _, argument := range call.Arguments {
				if i > 0 {
					builder.WriteString(" ")
				}
				switch argument.(type) {
				case *goja.Object:
					var marshalData []byte
					object := argument.ToObject(runtime)
					marshalData, err = object.MarshalJSON()
					if err != nil {
						config.ErrorLog.Println(err, stack.Trace())
						return goja.Undefined()
					}
					builder.WriteString(string(marshalData))
				default:
					value := argument.String()
					if value == "https://svelte.dev/e/experimental_async_ssr" {
						// Skipping experimental async ssr warnings.
						return goja.Undefined()
					}
					builder.WriteString(value)
				}
				i++
			}
			logger.Println(builder.String())
			return goja.Undefined()
		}
	}

	if err = console.Set("log", createLogger(LogLevelBase)); err != nil {
		return
	}

	if err = console.Set("info", createLogger(LogLevelBase)); err != nil {
		return
	}

	if err = console.Set("warn", createLogger(LogLevelWarning)); err != nil {
		return
	}

	if err = console.Set("error", createLogger(LogLevelDanger)); err != nil {
		return
	}

	if err = runtime.Set("console", console); err != nil {
		return
	}

	var text string
	if text, err = js.Bundle(filepath.Join(config.App, "dist"), api.FormatCommonJS, string(config.Data)); err != nil {
		return
	}

	var prog *goja.Program
	if prog, err = goja.Compile(server, fmt.Sprintf(config.Format, text), false); err != nil {
		return
	}

	var value goja.Value
	if value, err = runtime.RunProgram(prog); err != nil {
		return
	}

	var isfun bool
	if render, isfun = goja.AssertFunction(value); !isfun {
		err = errors.New("render is not a function")
	}

	return
}
