//go:build trace

package stacks

import (
	"fmt"
	"runtime"
	"strings"
)

var TraceSize = 10

// Trace returns the stack trace including the file name and line number.
func Trace() string {
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
