package js

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"rogchap.com/v8go"
)

var cache = map[string]*v8go.CompilerCachedData{}

// JavaScriptInvalidate invalidates compiler cached data.
//
// Whenever you invoke JavaScriptRun, cache data is extracted from the script,
// which is then used to speed up execution the next time a given script is executed.
func JavaScriptInvalidate(id string) {
	delete(cache, id)
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
func JavaScriptRun(id string, sourceCode []byte, globalFunctions map[string]v8go.FunctionCallback) (*v8go.Value, func(), error) {
	isolate := v8go.NewIsolate()
	globals := v8go.NewObjectTemplate(isolate)

	for functionName, functionCallback := range globalFunctions {
		setError := globals.Set(functionName, v8go.NewFunctionTemplate(isolate, functionCallback))
		if setError != nil {
			return nil, nil, setError
		}
	}

	context := v8go.NewContext(isolate, globals)

	var script *v8go.UnboundScript

	data, exists := cache[id]
	if exists {
		compiledScript, compileError := isolate.CompileUnboundScript(string(sourceCode), id, v8go.CompileOptions{CachedData: data})
		if compileError != nil {
			return nil, nil, compileError
		}
		script = compiledScript
	} else {
		compiledScript, compileError := isolate.CompileUnboundScript(string(sourceCode), id, v8go.CompileOptions{})
		if compileError != nil {
			return nil, nil, compileError
		}
		cache[id] = compiledScript.CreateCodeCache()
		script = compiledScript
	}

	result, jsError := script.Run(context)
	if jsError != nil {
		return nil, func() {}, jsError
	}

	return result, func() {
		context.Close()
		isolate.Dispose()
	}, nil
}

// JavaScriptBundle bundles JavaScript source code into a specific format.
func JavaScriptBundle(root string, format api.Format, sourceCode []byte) (bundle []byte, err error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   string(sourceCode),
			ResolveDir: root,
		},
	})

	for _, buildError := range result.Errors {
		return nil, fmt.Errorf("%s in %s:%d:%d", buildError.Text, buildError.Location.File, buildError.Location.Line, buildError.Location.Column)
	}

	return result.OutputFiles[0].Contents, nil
}
