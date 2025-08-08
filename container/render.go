package container

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/js"
	"os"
	"strings"
)

// Render renders the application with the given properties.
func Render(a *Container, p map[string]any) (string, string, error) {
	var runtime *goja.Runtime
	var program *goja.Program
	var compileError error

	if a.Config.Development {
		runtime = goja.New()
		var fileNameFixed = strings.ReplaceAll(a.Config.Script, "\\", "/")
		data, rer := os.ReadFile(fileNameFixed)
		if rer != nil {
			return "", "", rer
		}
		source, bundleError := js.Bundle(a.Config.Root, api.FormatCommonJS, string(data))
		if bundleError != nil {
			return "", "", bundleError
		}
		program, compileError = goja.Compile(a.Config.Script, fmt.Sprintf(globals.RenderScriptFormat, source), false)
		if compileError != nil {
			return "", "", compileError
		}
	} else {
		runtime = <-a.Channels.Runtime
		program = <-a.Channels.Program
		defer func() { go func() { a.Channels.Runtime <- runtime }() }()
		defer func() { go func() { a.Channels.Program <- program }() }()
	}

	programResult, programError := runtime.RunProgram(program)
	if programError != nil {
		return "", "", programError
	}

	render, isFunction := goja.AssertFunction(programResult)

	if !isFunction {
		return "", "", errors.New("render is not a function")
	}

	promise, programError := render(goja.Undefined(), runtime.ToValue(p))

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
