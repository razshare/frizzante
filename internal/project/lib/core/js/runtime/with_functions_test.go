package runtime

import (
	"main/lib/core/js"
	"testing"

	"github.com/dop251/goja"
)

func TestWithFunctions(t *testing.T) {
	var invoked1 bool
	var invoked2 bool
	run := goja.New()

	err := WithFunctions(run, map[string]js.Function{
		"custom_function_1": func(call goja.FunctionCall) goja.Value {
			invoked1 = true
			return goja.Undefined()
		},
		"custom_function_2": func(call goja.FunctionCall) goja.Value {
			invoked2 = true
			return goja.Undefined()
		},
	})
	if err != nil {
		return
	}

	_, err = run.RunString("custom_function_1();custom_function_2()")
	if err != nil {
		t.Fatal(err)
	}

	if !invoked1 {
		t.Fatal("custom_function_1 should be invoked")
	}

	if !invoked2 {
		t.Fatal("custom_function_2 should be invoked")
	}
}

func TestWithFunctionsShouldFail(t *testing.T) {
	run := goja.New()
	err := WithFunctions(run, map[string]js.Function{
		"?''^_@_-custom_function_1": func(call goja.FunctionCall) goja.Value {
			return goja.Undefined()
		},
		"custom_function_2": func(call goja.FunctionCall) goja.Value {
			return goja.Undefined()
		},
	})
	if err != nil {
		t.Fatal("functions should fail assignment")
		return
	}
}
