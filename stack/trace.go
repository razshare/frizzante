package stack

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var TraceEnabled = os.Getenv("DEV") == "1" || os.Getenv("TRACE") == "1"
var TraceSize = 10

// Trace returns the stack trace including the file name and line number.
func Trace() string {
	if !TraceEnabled {
		return ""
	}

	var sb strings.Builder
	ptr := make([]uintptr, TraceSize)
	runtime.Callers(2, ptr)
	frames := runtime.CallersFrames(ptr)

	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return sb.String()
}
