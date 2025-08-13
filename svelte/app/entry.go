package app

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

// RunEntry runs the container entry point with the given properties.
//
// The container's entry point is usually Svelte's render function.
func RunEntry(c *Config, p map[string]any) (string, string, error) {
	var rt *goja.Runtime
	var prg *goja.Program
	var cerr error

	if c.Development {
		rt = goja.New()

		nfix := strings.ReplaceAll(c.Script, "\\", "/")

		data, rer := os.ReadFile(nfix)
		if rer != nil {
			return "", "", rer
		}

		src, berr := js.Bundle(c.Root, api.FormatCommonJS, string(data))
		if berr != nil {
			return "", "", berr
		}

		prg, cerr = goja.Compile(c.Script, fmt.Sprintf(globals.RenderScriptFormat, src), false)
		if cerr != nil {
			return "", "", cerr
		}
	} else {
		rt = <-c.Channels.Runtime
		prg = <-c.Channels.Program
		defer func() { go func() { c.Channels.Runtime <- rt }() }()
		defer func() { go func() { c.Channels.Program <- prg }() }()
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
