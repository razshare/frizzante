package messages

import (
	"fmt"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/tui/config"
	"os"
	"strings"
)

func Fatal(args ...any) {
	l := len(args)
	entries := make([]string, l+1)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	entries[l] = "\n" + stack.Trace()
	Status("ERROR", strings.Join(entries, ""), config.Colors.Error, "15", config.Colors.Error)
	os.Exit(1)
}

func Fatalf(format string, vars ...any) {
	Fatal(fmt.Sprintf(format, vars...), "\n", stack.Trace())
}
