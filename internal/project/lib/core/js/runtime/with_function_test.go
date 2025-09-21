package runtime

import (
	"testing"

	"github.com/dop251/goja"
)

func TestWithFunction(t *testing.T) {
	var err error
	var invoked bool

	runtime := goja.New()

	if err = WithFunction(runtime, "custom_function", func(call goja.FunctionCall) goja.Value {
		invoked = true
		return goja.Undefined()
	}); err != nil {
		return
	}

	if _, err = runtime.RunString("custom_function()"); err != nil {
		t.Fatal(err)
	}

	if !invoked {
		t.Fatal("custom_function should be invoked")
	}
}
