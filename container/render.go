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
	var rt *goja.Runtime
	var prg *goja.Program
	var cerr error

	if a.Config.Development {
		rt = goja.New()

		nfix := strings.ReplaceAll(a.Config.Script, "\\", "/")

		data, rer := os.ReadFile(nfix)
		if rer != nil {
			return "", "", rer
		}

		src, berr := js.Bundle(a.Config.Root, api.FormatCommonJS, string(data))
		if berr != nil {
			return "", "", berr
		}

		prg, cerr = goja.Compile(a.Config.Script, fmt.Sprintf(globals.RenderScriptFormat, src), false)
		if cerr != nil {
			return "", "", cerr
		}
	} else {
		rt = <-a.Channels.Runtime
		prg = <-a.Channels.Program
		defer func() { go func() { a.Channels.Runtime <- rt }() }()
		defer func() { go func() { a.Channels.Program <- prg }() }()
	}

	pres, perr := rt.RunProgram(prg)
	if perr != nil {
		return "", "", perr
	}

	render, isf := goja.AssertFunction(pres)

	if !isf {
		return "", "", errors.New("render is not a function")
	}

	promise, perr := render(goja.Undefined(), rt.ToValue(p))

	if perr != nil {
		return "", "", perr
	}

	value := promise.Export().(*goja.Promise).Result().ToObject(rt)

	headv := value.Get("head")
	bodyv := value.Get("body")

	var head string
	var body string

	if headv != nil {
		head = headv.String()
	}

	if bodyv != nil {
		body = bodyv.String()
	}

	return head, body, nil
}
