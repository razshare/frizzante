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

// JavaScriptContextCreateWithGlobals creates a JavaScript context with a map of global functions.
func JavaScriptContextCreateWithGlobals(globals map[string]v8go.FunctionCallback) (*JavaScriptContext, error) {
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

var v8MapOfCompilerCachedData = map[string]*v8go.CompilerCachedData{}

// JavaScriptInvalidateCompilerCacheData invalidates compiler code cached data.
//
// Whenever you invoke JavaScriptRun, cache data is extracted from the script,
// which is then used to speed up execution the next time a given script is executed.
func JavaScriptInvalidateCompilerCacheData(id string) {
	delete(v8MapOfCompilerCachedData, id)
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
func JavaScriptRun(id string, source string, globals map[string]v8go.FunctionCallback) (
	result *v8go.Value,
	destroy func(),
	scriptError error,
) {
	js, createError := JavaScriptContextCreateWithGlobals(globals)
	if createError != nil {
		return nil, func() {}, createError
	}

	var script *v8go.UnboundScript

	codeCache, hasCodeCash := v8MapOfCompilerCachedData[id]
	if hasCodeCash {
		compiledScript, compilationError := js.isolate.CompileUnboundScript(source, id, v8go.CompileOptions{CachedData: codeCache})
		if compilationError != nil {
			return nil, nil, compilationError
		}
		script = compiledScript
	} else {
		compiledScript, compilationError := js.isolate.CompileUnboundScript(source, id, v8go.CompileOptions{})
		if compilationError != nil {
			return nil, nil, compilationError
		}
		v8MapOfCompilerCachedData[id] = compiledScript.CreateCodeCache()
		script = compiledScript
	}

	scriptResult, scriptError := script.Run(js.context)
	if scriptError != nil {
		return nil, func() {}, scriptError
	}

	return scriptResult, func() { JavaScriptDestroy(js) }, nil
}

// JavaScriptDestroy destroys a JavaScriptContext.
//
// Usually you don't need to invoke this manually because JavaScriptRun returns a destroyer function.
func JavaScriptDestroy(js *JavaScriptContext) {
	js.context.Close()
	js.isolate.Dispose()
}

// JavaScriptBundle bundles JavaScript source code into a specific format.
func JavaScriptBundle(rootDirectory string, format api.Format, source string) (bundle string, bundleError error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   source,
			ResolveDir: rootDirectory,
		},
	})

	for _, err := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", err.Text, err.Location.File, err.Location.Line, err.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
