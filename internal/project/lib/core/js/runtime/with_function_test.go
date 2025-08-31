package runtime

import (
	"testing"

	"github.com/dop251/goja"
)

func TestWithFunction(t *testing.T) {
	var invoked bool
	run := goja.New()
	err := WithFunction(run, "custom_function", func(call goja.FunctionCall) goja.Value {
		invoked = true
		return goja.Undefined()
	})
	if err != nil {
		return
	}

	_, err = run.RunString("custom_function()")
	if err != nil {
		t.Fatal(err)
	}

	if !invoked {
		t.Fatal("custom_function should be invoked")
	}
}
