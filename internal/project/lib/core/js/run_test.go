package js

import (
	"slices"
	"strings"
	"testing"

	"github.com/dop251/goja"
)

func TestRun(t *testing.T) {
	run := goja.New()
	src := "1+1"
	val, err := run.RunString(src)
	if err != nil {
		t.Fatal(err)
	}

	if val.ToInteger() != 2 {
		t.Fatal("value should be 2")
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
	val, err = run.RunString(src)
	if err != nil {
		t.Fatal(err)
	}

	obj := val.ToObject(run)
	keys := obj.Keys()

	if !slices.Contains(keys, "long") {
		t.Fatal("value should have a 'long' key")
	}

	if !slices.Contains(keys, "short") {
		t.Fatal("value should have a 'short' key")
	}

	long := obj.Get("long")
	short := obj.Get("short")

	longs := strings.Split(long.String(), "-")
	if len(longs) != 5 {
		t.Fatal("long string should composed of 5 part separated by 4 -")
	}

	shorts := strings.Split(short.String(), "-")
	if len(shorts) != 1 {
		t.Fatal("string should not contain any -", short.String())
	}

}
