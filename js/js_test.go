package js

import (
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestJavaScriptRun(t *testing.T) {
	rt := goja.New()
	src := "1+1"
	ac, err := rt.RunString(src)
	if err != nil {
		t.Fatal(err)
	}

	if ac.ToInteger() != 2 {
		t.Fatalf("script was expected to return 2, received '%d' instead", ac.ToInteger())
	}

	src = `
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
	ac, err = rt.RunString(src)
	if err != nil {
		t.Fatal(err)
	}

	obj := ac.ToObject(rt)
	keys := obj.Keys()

	if !slices.Contains(keys, "long") {
		t.Fatal("actual value was expected to have a 'long' key")
	}

	if !slices.Contains(keys, "short") {
		t.Fatal("actual value was expected to have a 'short' key")
	}

	long := obj.Get("long")
	short := obj.Get("short")

	longs := strings.Split(long.String(), "-")
	if len(longs) != 5 {
		t.Fatalf("long string was expected to be composed of 5 part separated by 4 -, received '%s' instead", long.String())
	}

	shorts := strings.Split(short.String(), "-")
	if len(shorts) != 1 {
		t.Fatalf("string was expected to be composed of 1 part, received '%s' instead", short.String())
	}

}

func TestJavaScriptBundle(t *testing.T) {
	ac := ""
	ex := "hello"

	rt := goja.New()

	err := SetFunction(rt, "signal", func(call goja.FunctionCall) goja.Value {
		args := call.Arguments
		if len(args) > 0 {
			ac = args[0].String()
		}
		return goja.Undefined()
	})

	if err != nil {
		return
	}

	src := `
	import { writable } from 'svelte/store'
	const test = writable("hello")
	test.subscribe(function updated(value){
		signal(value)
	})
	`

	cjs, err := Bundle(filepath.Join("app"), api.FormatCommonJS, src)
	if err != nil {
		t.Fatal(err)
	}

	_, err = rt.RunString(cjs)
	if err != nil {
		t.Fatal(err)
	}

	if ac != ex {
		t.Fatalf("script was expected to update the actual value to '%s', received '%s' instead.", ex, ac)
	}
}
