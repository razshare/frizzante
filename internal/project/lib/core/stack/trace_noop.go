//go:build !trace

package stack

var TraceSize = 10

// Trace returns the stack trace including the file name and line number.
func Trace() string {
	return ""
}
