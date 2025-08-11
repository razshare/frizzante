package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Tip(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("TIP", strings.Join(entries, ""), config.Colors.Tip, "15", config.Colors.Tip)
}

func Tipf(format string, vars ...any) {
	Info(fmt.Sprintf(format, vars...))
}
