package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Info(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("INFO", strings.Join(entries, ""), config.Colors.Info, "15", config.Colors.Info)
}

func Infof(format string, vars ...any) {
	Info(fmt.Sprintf(format, vars...))
}
