package frizzante

import (
	"github.com/evanw/esbuild/pkg/api"
	"rogchap.com/v8go"
	"strings"
	"testing"
)

func TestNewJavaScriptContext(test *testing.T) {
	_, err := newJavaScriptContext(map[string]v8go.FunctionCallback{})
	if err != nil {
		test.Fatal(err)
	}
}

func TestJavaScriptRun(test *testing.T) {
	// Simple.
	script := "1+1"
	actual, destroy, javaScriptError := JavaScriptRun(script, map[string]v8go.FunctionCallback{})
	defer destroy()
	if javaScriptError != nil {
		test.Fatal(javaScriptError)
	}

	if actual.Int32() != 2 {
		test.Fatalf("script was expected to return 2, received '%d' instead", actual.Int32())
	}

	// Complex and with JsDoc.
	script = `
	/**
	 * @param {boolean} payload
	 * @returns
	 */
	function uuid(short = false) {
		let dt = new Date().getTime()
		const BLUEPRINT = short ? 'xyxxyxyx' : 'xxxxxxxx-xxxx-yxxx-yxxx-xxxxxxxxxxxx'
		const RESULT = BLUEPRINT.replace(/[xy]/g, function run(c) {
		const r = (dt + Math.random() * 16) % 16 | 0
		dt = Math.floor(dt / 16)
		return (c == 'x' ? r : (r & 0x3) | 0x8).toString(16)
		})
		return RESULT
	}
	
	const result = {
		long: uuid(),
		short: uuid(true),
	}
	
	result
	`
	actual, destroy, javaScriptError = JavaScriptRun(script, map[string]v8go.FunctionCallback{})
	defer destroy()
	if javaScriptError != nil {
		test.Fatal(javaScriptError)
	}

	obj := actual.Object()

	if !obj.Has("long") {
		test.Fatal("actual value was expected to have a 'long' key")
	}

	if !obj.Has("short") {
		test.Fatal("actual value was expected to have a 'short' key")
	}

	long, longError := obj.Get("long")
	if longError != nil {
		test.Fatal(longError)
	}
	short, shortError := obj.Get("short")
	if shortError != nil {
		test.Fatal(shortError)
	}

	longPieces := strings.Split(long.String(), "-")
	if len(longPieces) != 5 {
		test.Fatalf("long string was expected to be composed of 5 part separated by 4 -, received '%s' instead", long.String())
	}

	shortPieces := strings.Split(short.String(), "-")
	if len(shortPieces) != 1 {
		test.Fatalf("string was expected to be composed of 1 part, received '%s' instead", short.String())
	}

}

func TestJavaScriptBundle(test *testing.T) {
	script := `
	import { writable } from 'svelte/store'
	const test = writable("hello")
	test.subscribe(function updated(value){
		signal(value)
	})
	`

	cjs, bundleError := JavaScriptBundle(".", api.FormatCommonJS, script)
	if bundleError != nil {
		test.Fatal(bundleError)
	}
	actual := ""
	expected := "hello"
	_, destroy, javaScriptError := JavaScriptRun(cjs, map[string]v8go.FunctionCallback{
		"signal": func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			args := info.Args()
			if len(args) > 0 {
				actual = args[0].String()
			}
			return nil
		},
	})
	defer destroy()
	if javaScriptError != nil {
		test.Fatal(javaScriptError)
	}

	if actual != expected {
		test.Fatalf("script was expected to update the actual value to '%s', received '%s' instead.", expected, actual)
	}
}
