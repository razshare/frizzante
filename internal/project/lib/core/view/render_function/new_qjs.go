//go:build experimental_qjs_runtime

package render_function

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fastschema/qjs"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

func New(config Config) (render RenderFunction, err error) {
	var runtime *qjs.Runtime
	var builder strings.Builder

	if runtime, err = qjs.New(); err != nil {
		return
	}

	context := runtime.Context()

	createLogger := func(level LogLevel) *qjs.Value {
		var logger *log.Logger

		switch level {
		case LogLevelDanger:
			logger = config.ErrorLog
		default:
			logger = config.InfoLog
		}

		return context.Function(
			func(this *qjs.This) (_ *qjs.Value, _ error) {
				builder.Reset()
				i := 0
				for _, argument := range this.Args() {
					if i > 0 {
						builder.WriteString(" ")
					}
					if argument.IsObject() {
						var marshalData string

						marshalData, err = argument.JSONStringify()
						if err != nil {
							config.ErrorLog.Println(err, stack.Trace())
							return
						}
						builder.WriteString(marshalData)
						continue
					}

					value := argument.String()
					if value == "https://svelte.dev/e/experimental_async_ssr" {
						// Skipping experimental async ssr warnings.
						return
					}
					builder.WriteString(value)

					i++
				}
				logger.Println(builder.String())
				return
			},
			false,
		)
	}

	console := context.NewObject()
	console.SetProperty(context.NewString("log"), createLogger(LogLevelBase))
	console.SetProperty(context.NewString("info"), createLogger(LogLevelBase))
	console.SetProperty(context.NewString("warn"), createLogger(LogLevelWarning))
	console.SetProperty(context.NewString("error"), createLogger(LogLevelDanger))

	context.Global().SetProperty(context.NewString("console"), console)

	var renderValue *qjs.Value
	context.SetFunc("frizzante_set_render", func(this *qjs.This) (*qjs.Value, error) {
		renderValue = this.Args()[0]
		return nil, nil
	})

	var text string
	if text, err = js.Bundle(filepath.Join(config.App, "dist"), api.FormatCommonJS, string(config.Data)); err != nil {
		return
	}

	source := "const module={exports:{}};\n" + text + "\nfrizzante_set_render(render)"

	if _, err = context.Eval("app.server.js", qjs.Code(source)); err != nil {
		return
	}

	if !renderValue.IsFunction() {
		err = errors.New("render is not a function")
		return
	}

	render = func(view view_.View) (head string, body string, err error) {
		var propsObject *qjs.Value
		if propsObject, err = qjs.ToJSValue(context, view_.NewData(view)); err != nil {
			return
		}

		var promise *qjs.Value
		if promise, err = context.Invoke(renderValue, context.Global(), propsObject); err != nil {
			return
		}

		var value *qjs.Value
		if value, err = promise.Await(); err != nil {
			return
		}
		defer value.Free()

		headv := value.GetProperty(context.NewString("head"))
		bodyv := value.GetProperty(context.NewString("body"))

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
