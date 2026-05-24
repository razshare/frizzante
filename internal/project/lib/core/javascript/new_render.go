package javascript

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/security"
	"github.com/razshare/frizzante/internal/project/lib/core/types"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func NewRender(options NewRenderOptions) (render Render, err error) {
	runtime := goja.New()
	console := runtime.NewObject()
	if err = console.Set("log", CreateLogger(CreateLoggerOptions{
		Level:    LogLevelBase,
		Runtime:  runtime,
		ErrorLog: options.ErrorLog,
		InfoLog:  options.InfoLog,
	})); err != nil {
		return
	}
	if err = console.Set("info", CreateLogger(CreateLoggerOptions{
		Level:    LogLevelBase,
		Runtime:  runtime,
		ErrorLog: options.ErrorLog,
		InfoLog:  options.InfoLog,
	})); err != nil {
		return
	}
	if err = console.Set("warn", CreateLogger(CreateLoggerOptions{
		Level:    LogLevelWarning,
		Runtime:  runtime,
		ErrorLog: options.ErrorLog,
		InfoLog:  options.InfoLog,
	})); err != nil {
		return
	}
	if err = console.Set("error", CreateLogger(CreateLoggerOptions{
		Level:    LogLevelDanger,
		Runtime:  runtime,
		ErrorLog: options.ErrorLog,
		InfoLog:  options.InfoLog,
	})); err != nil {
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
	if err = runtime.Set("hash_256_627f0c8d3f2158f776f550ab47b35de9", func(call goja.FunctionCall) goja.Value {
		text := call.Arguments[1].String()
		return runtime.ToValue(security.Sha256(text))
	}); err != nil {
		return
	}
	var source string
	if source, err = options.FindSource(); err != nil {
		return
	}
	const bootstrap = `
		const module={exports:{}};
		const globalThis={crypto:{subtle:{digest:hash_256_627f0c8d3f2158f776f550ab47b35de9}}};
	`
	script := fmt.Sprintf("%s\n%s\nfrizzante_set_render(render)", strings.TrimSpace(bootstrap), source)
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
