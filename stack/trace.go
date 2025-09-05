package stack

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var TraceSize = 10
var TraceEnabled = os.Getenv("DEV") == "1" || os.Getenv("TRACE") == "1"

// Trace returns the stack trace including the file name and line number.
func Trace() string {
	if !TraceEnabled {
		return ""
	}

	var builder strings.Builder
	ptrs := make([]uintptr, TraceSize)
	runtime.Callers(2, ptrs)
	frames := runtime.CallersFrames(ptrs)

	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return builder.String()
}
