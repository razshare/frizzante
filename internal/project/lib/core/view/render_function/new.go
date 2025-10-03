//go:build !experimental_qjs_runtime

package render_function

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	view_ "github.com/razshare/frizzante/internal/project/lib/core/view"
)

func New(config Config) (render RenderFunction, err error) {
	var builder strings.Builder

	runtime := goja.New()
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

	var renderValue goja.Value
	if err = runtime.Set("frizzante_set_render", func(call goja.FunctionCall) goja.Value {
		renderValue = call.Arguments[0]
		return goja.Undefined()
	}); err != nil {
		return
	}

	var text string
	if text, err = js.Bundle(filepath.Join(config.App, "dist"), api.FormatCommonJS, string(config.Data)); err != nil {
		return
	}

	source := "const module={exports:{}};\n" + text + "\nfrizzante_set_render(render)"

	var prog *goja.Program
	if prog, err = goja.Compile("app.server.js", source, false); err != nil {
		return
	}

	if _, err = runtime.RunProgram(prog); err != nil {
		return
	}

	var isfun bool
	var renderJs goja.Callable
	if renderJs, isfun = goja.AssertFunction(renderValue); !isfun {
		err = errors.New("render is not a function")
		return
	}

	render = func(view view_.View) (head string, body string, err error) {
		var props map[string]any
		if props, err = types.EncodeInterface(view_.NewData(view)); err != nil {
			return
		}

		var promise goja.Value
		if promise, err = renderJs(goja.Undefined(), runtime.ToValue(props)); err != nil {
			return
		}

		result := promise.Export().(*goja.Promise).Result().ToObject(runtime)

		headv := result.Get("head")
		bodyv := result.Get("body")

		if headv != nil {
			head = headv.String()
		}

		if bodyv != nil {
			body = bodyv.String()
		}

		return
	}

	return
}
