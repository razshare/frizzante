package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/config"
)

func Info(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("INFO", strings.Join(entries, ""), config.Colors.Info, "17", config.Colors.Info)
}

func Infof(format string, vars ...any) {
	Info(fmt.Sprintf(format, vars...))
}
