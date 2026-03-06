package javascript

import (
	"log"
	"strings"

	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func CreateLogger(options CreateLoggerOptions) func(call goja.FunctionCall) goja.Value {
	var builder strings.Builder
	var logger *log.Logger
	switch options.Level {
	case LogLevelDanger:
		logger = options.ErrorLog
	default:
		logger = options.InfoLog
	}
	return func(call goja.FunctionCall) goja.Value {
		builder.Reset()
		i := 0
		for _, argument := range call.Arguments {
			if i > 0 {
				builder.WriteString(" ")
			}
			switch argument.(type) {
			case *goja.Object:
				var marshalData []byte
				object := argument.ToObject(options.Runtime)
				var err error
				marshalData, err = object.MarshalJSON()
				if err != nil {
					options.ErrorLog.Printf(
					"javascript.CreateLogger: failed to marshal JavaScript object to JSON: %v\n%s",
					err,
					stack.Trace(),
				)
					return goja.Undefined()
				}
				builder.WriteString(string(marshalData))
			default:
				value := argument.String()
				if value == "https://svelte.dev/e/experimental_async_ssr" {
					// Skipping experimental async ssr warnings.
					return goja.Undefined()
				}
				builder.WriteString(value)
			}
			i++
		}
		logger.Println(builder.String())
		return goja.Undefined()
	}
}
