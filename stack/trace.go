package stack

import (
	"fmt"
	"runtime"
)

// Trace returns the stack trace including the file name and line number.
func Trace() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return ""
}
