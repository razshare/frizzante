package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Warning(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("WARNING", strings.Join(entries, ""), config.Colors.Warning, "0", config.Colors.Warning)
}

func Warningf(format string, vars ...any) {
	Warning(fmt.Sprintf(format, vars...))
}
