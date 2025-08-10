package main

import (
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/js"
	"slices"
	"strings"
	"testing"
)

func TestJavaScriptRun(test *testing.T) {
	runtime := goja.New()
	script := "1+1"
	actual, runError := runtime.RunString(script)
	if runError != nil {
		test.Fatal(runError)
	}

	if actual.ToInteger() != 2 {
		test.Fatalf("script was expected to return 2, received '%d' instead", actual.ToInteger())
	}

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
	actual, runError = runtime.RunString(script)
	if runError != nil {
		test.Fatal(runError)
	}

	obj := actual.ToObject(runtime)
	keys := obj.Keys()

	if !slices.Contains(keys, "long") {
		test.Fatal("actual value was expected to have a 'long' key")
	}

	if !slices.Contains(keys, "short") {
		test.Fatal("actual value was expected to have a 'short' key")
	}

	long := obj.Get("long")
	short := obj.Get("short")

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
	actual := ""
	expected := "hello"

	runtime := goja.New()

	err := js.SetFunction(runtime, "signal", func(call goja.FunctionCall) goja.Value {
		args := call.Arguments
		if len(args) > 0 {
			actual = args[0].String()
		}
		return goja.Undefined()
	})
	if err != nil {
		return
	}

	script := `
	import { writable } from 'svelte/store'
	const test = writable("hello")
	test.subscribe(function updated(value){
		signal(value)
	})
	`

	cjs, bundleError := js.Bundle("app", api.FormatCommonJS, script)
	if bundleError != nil {
		test.Fatal(bundleError)
	}

	_, javaScriptError := runtime.RunString(cjs)
	if javaScriptError != nil {
		test.Fatal(javaScriptError)
	}

	if actual != expected {
		test.Fatalf("script was expected to update the actual value to '%s', received '%s' instead.", expected, actual)
	}
}
