package frizzante

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"rogchap.com/v8go"
)

type JavaScriptContext struct {
	isolate *v8go.Isolate
	global  *v8go.ObjectTemplate
	context *v8go.Context
}

func newJavaScriptContext(globals map[string]v8go.FunctionCallback) (*JavaScriptContext, error) {
	isolate := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(isolate)

	for key, callback := range globals {
		template := v8go.NewFunctionTemplate(isolate, callback)
		setError := global.Set(key, template)
		if setError != nil {
			return nil, setError
		}
	}

	context := v8go.NewContext(isolate, global)

	return &JavaScriptContext{
		isolate: isolate,
		global:  global,
		context: context,
	}, nil
}

// JavaScriptRun runs a javascript module.
//
// It returns the last expression of the script and a destroyer function.
//
// The destroyer function is never nil.
//
// You should always call the destroyer function as soon as possible to limit memory usage.
//
// Each global function will be injected into the context of the module automatically so that you can invoke them from the script.
func JavaScriptRun(source string, globals map[string]v8go.FunctionCallback) (
	result *v8go.Value,
	destroy func(),
	runError error,
) {
	js, createError := newJavaScriptContext(globals)
	if createError != nil {
		return nil, func() {}, createError
	}
	exports, runError := js.context.RunScript(source, "frizzante.js")
	if runError != nil {
		return nil, func() {}, runError
	}

	return exports, func() { JavaScriptDestroy(js) }, nil
}

func JavaScriptDestroy(js *JavaScriptContext) {
	js.context.Close()
	js.isolate.Dispose()
}

func JavaScriptBundle(root string, format api.Format, source string) (bundle string, bundleError error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   source,
			ResolveDir: root,
		},
	})

	for _, err := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", err.Text, err.Location.File, err.Location.Line, err.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
