package js

import (
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"path/filepath"
	"testing"
)

func TestBundle(t *testing.T) {
	run := goja.New()

	src := `
	import { writable } from 'svelte/store'
	const test = writable("hello")
	let result
	function done(){
		return result
	}
	test.subscribe(function updated(value){
		result = value
	})
	done()
	`

	cjs, err := Bundle(filepath.Join("app"), api.FormatCommonJS, src)
	if err != nil {
		t.Fatal(err)
	}

	var val goja.Value
	val, err = run.RunString(cjs)
	if err != nil {
		t.Fatal(err)
	}

	if val.String() != "hello" {
		t.Fatalf("value should be hello")
	}
}

func TestBundleShouldFail(t *testing.T) {
	src := `
	import { writable } from 'svelte/store___' //invalid import, this should fail
	const test = writable("hello")
	let result
	function done(){
		return result
	}
	test.subscribe(function updated(value){
		result = value
	})
	done()
	`

	_, err := Bundle(filepath.Join("app"), api.FormatCommonJS, src)
	if err == nil {
		t.Fatal("bundle should fail")
	}
}
