package messages

import (
	"fmt"
	"github.com/razshare/frizzante/stack"
	"os"
)

func Fatal(args ...any) {
	Error(args...)
	os.Exit(1)
}

func Fatalf(format string, vars ...any) {
	Fatal(fmt.Sprintf(format, vars...), "\n", stack.Trace())
}
