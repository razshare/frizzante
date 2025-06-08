package frz

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"rogchap.com/v8go"
)

var javaScriptCache = map[string]*v8go.CompilerCachedData{}

// JavaScriptInvalidate invalidates compiler cached data.
//
// Whenever you invoke JavaScriptRun, cache data is extracted from the script,
// which is then used to speed up execution the next time a given script is executed.
func JavaScriptInvalidate(id string) {
	delete(javaScriptCache, id)
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
func JavaScriptRun(id string, source []byte, functions map[string]v8go.FunctionCallback) (
	result *v8go.Value,
	destroy func(),
	scriptError error,
) {
	isolate := v8go.NewIsolate()
	globals := v8go.NewObjectTemplate(isolate)

	for key, callback := range functions {
		setError := globals.Set(key, v8go.NewFunctionTemplate(isolate, callback))
		if setError != nil {
			return nil, nil, setError
		}
	}

	context := v8go.NewContext(isolate, globals)

	var script *v8go.UnboundScript

	codeCache, hasCodeCash := javaScriptCache[id]
	if hasCodeCash {
		compiledScript, compilationError := isolate.CompileUnboundScript(string(source), id, v8go.CompileOptions{CachedData: codeCache})
		if compilationError != nil {
			return nil, nil, compilationError
		}
		script = compiledScript
	} else {
		compiledScript, compilationError := isolate.CompileUnboundScript(string(source), id, v8go.CompileOptions{})
		if compilationError != nil {
			return nil, nil, compilationError
		}
		javaScriptCache[id] = compiledScript.CreateCodeCache()
		script = compiledScript
	}

	scriptResult, scriptError := script.Run(context)
	if scriptError != nil {
		return nil, func() {}, scriptError
	}

	return scriptResult, func() {
		context.Close()
		isolate.Dispose()
	}, nil
}

// JavaScriptBundle bundles JavaScript source code into a specific format.
func JavaScriptBundle(rootDirectory string, format api.Format, source []byte) (bundle []byte, bundleError error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   string(source),
			ResolveDir: rootDirectory,
		},
	})

	for _, err := range result.Errors {
		return nil, fmt.Errorf("%s in %s:%d:%d", err.Text, err.Location.File, err.Location.Line, err.Location.Column)
	}

	return result.OutputFiles[0].Contents, nil
}
