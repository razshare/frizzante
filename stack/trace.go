package stack

import (
	"fmt"
	"runtime"
)

// Trace returns the stack trace including the file name and line number.
func Trace() string {
	if _, file, line, ok := runtime.Caller(2); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return ""
}
