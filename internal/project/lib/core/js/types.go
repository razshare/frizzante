package js

import "github.com/dop251/goja"

type Function = func(call goja.FunctionCall) goja.Value
