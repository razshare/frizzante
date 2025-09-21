package runtime

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
)

func TestWithFunctions(t *testing.T) {
	var err error
	var invoked1 bool
	var invoked2 bool
	runtime := goja.New()

	if err = WithFunctions(runtime, map[string]js.Function{
		"custom_function_1": func(call goja.FunctionCall) goja.Value {
			invoked1 = true
			return goja.Undefined()
		},
		"custom_function_2": func(call goja.FunctionCall) goja.Value {
			invoked2 = true
			return goja.Undefined()
		},
	}); err != nil {
		return
	}

	if _, err = runtime.RunString("custom_function_1();custom_function_2()"); err != nil {
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
	var err error

	runtime := goja.New()

	if err = WithFunctions(runtime, map[string]js.Function{
		"?''^_@_-custom_function_1": func(call goja.FunctionCall) goja.Value {
			return goja.Undefined()
		},
		"custom_function_2": func(call goja.FunctionCall) goja.Value {
			return goja.Undefined()
		},
	}); err != nil {
		t.Fatal("functions should fail assignment")
		return
	}
}
