package frizzante

import (
	"embed"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"os"
	"path/filepath"
	"rogchap.com/v8go"
)

type JavaScript struct {
	isolate *v8go.Isolate
	global  *v8go.ObjectTemplate
	context *v8go.Context
}

// NewJavaScript creates a JavaScript with a map of global functions.
func NewJavaScript(globals map[string]v8go.FunctionCallback) (*JavaScript, error) {
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

	return &JavaScript{
		isolate: isolate,
		global:  global,
		context: context,
	}, nil
}

// Destroy destroys the context.
//
// Usually you don't need to invoke this manually because JavaScriptRun returns a destroyer function.
func (js *JavaScript) Destroy() {
	js.context.Close()
	js.isolate.Dispose()
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
	js, createError := NewJavaScript(globals)
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

	return scriptResult, func() { js.Destroy() }, nil
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

// JavaScriptRender executes the `.dist/server/render.server.ts` file
// and returns the head of the document along with its body.
//
// If the environment variable DEV is set to 1, the file .dist/server/render.server.ts is executed directly from the
// local file system, otherwise RenderServerJs executes the file .dist/server/render.server.ts located within the
// view's embedded file system.
func JavaScriptRender(efs embed.FS, stringifiedProps string) (head string, body string, jsError error) {
	renderFileName := filepath.Join(".dist", "server", "render.server.js")

	var renderEsmBytes []byte
	if "1" == os.Getenv("DEV") {
		renderEsmBytesLocal, readError := os.ReadFile(renderFileName)
		if readError != nil {
			return "", "", readError
		}
		renderEsmBytes = renderEsmBytesLocal
	} else {
		renderEsmBytesLocal, readError := efs.ReadFile(renderFileName)
		if readError != nil {
			return "", "", readError
		}
		renderEsmBytes = renderEsmBytesLocal
	}

	renderEsm := string(renderEsmBytes)

	renderCjs, javaScriptBundleError := JavaScriptBundle(".", api.FormatCommonJS, renderEsm)
	if javaScriptBundleError != nil {
		return "", "", javaScriptBundleError
	}

	renderIif := fmt.Sprintf("const module={exports:{}}; const render = \n(function(){\n%s\nreturn render;\n})()", renderCjs)

	doneEsm := fmt.Sprintf(
		`
		%s
		render(JSON.parse(stringifiedProps())).then(function done(rendered){
			head(rendered.head??'');
			body(rendered.body??'');
		});
		`,
		renderIif,
	)

	doneCjs, bundleError := JavaScriptBundle(".", api.FormatCommonJS, doneEsm)
	if bundleError != nil {
		return "", "", bundleError
	}

	globals := map[string]v8go.FunctionCallback{}

	globals["stringifiedProps"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		value, valueError := v8go.NewValue(info.Context().Isolate(), stringifiedProps)
		if nil != valueError {
			return nil
		}
		return value
	}

	globals["inspect"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			message := args[0].String()
			println(message)
		}
		return nil
	}

	globals["head"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			head = args[0].String()
		}
		return nil
	}

	globals["body"] = func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			body = args[0].String()
		}
		return nil
	}

	_, destroy, javaScriptError := JavaScriptRun(renderFileName, doneCjs, globals)
	defer destroy()
	if javaScriptError != nil {
		return head, body, javaScriptError
	}

	return head, body, nil
}
