package javascript

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/dev/types"
)

func NewRender(options NewRenderOptions) (render Render, err error) {
	var builder strings.Builder
	runtime := goja.New()
	console := runtime.NewObject()
	createLogger := func(level LogLevel) func(call goja.FunctionCall) goja.Value {
		var logger *log.Logger
		switch level {
		case LogLevelDanger:
			logger = options.ErrorLog
		default:
			logger = options.InfoLog
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
						options.ErrorLog.Println(err, stack.Trace())
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
	var source string
	if source, err = options.FindSource(); err != nil {
		return
	}
	script := fmt.Sprintf("const module={exports:{}};\n%s\nfrizzante_set_render(render)", source)
	var prog *goja.Program
	if prog, err = goja.Compile("app.server.js", script, false); err != nil {
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
	render = func(options RenderOptions) (head string, body string, err error) {
		var props map[string]any
		if props, err = types.EncodeInterface(views.NewData(options.View)); err != nil {
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
