package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/config"
)

func Warning(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("WARNING", strings.Join(entries, ""), config.Colors.Warning, "0", config.Colors.Warning)
}

func Warningf(format string, vars ...any) {
	Warning(fmt.Sprintf(format, vars...))
}
