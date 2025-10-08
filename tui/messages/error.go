package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
	"github.com/razshare/frizzante/tui/config"
)

func Error(args ...any) {
	length := len(args)
	entries := make([]string, length+1)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}

	if trace := stacks.Trace(); trace != "" {
		entries[length] = "\n" + stacks.Trace()
	}

	Status("ERROR", strings.Join(entries, ""), config.Colors.Error, "233", config.Colors.Error)
}

func Errorf(format string, vars ...any) {
	Error(fmt.Sprintf(format, vars...))
}
