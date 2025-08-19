package messages

import (
	"fmt"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Error(args ...any) {
	l := len(args)
	entries := make([]string, l+1)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}

	if trace := stack.Trace(); trace != "" {
		entries[l] = "\n" + stack.Trace()
	}

	Status("ERROR", strings.Join(entries, ""), config.Colors.Error, "233", config.Colors.Error)
}

func Errorf(format string, vars ...any) {
	Error(fmt.Sprintf(format, vars...))
}
