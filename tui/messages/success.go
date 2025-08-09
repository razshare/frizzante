package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Success(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("SUCCESS", strings.Join(entries, ""), config.Colors.Success, "0", config.Colors.Success)
}

func Successf(format string, vars ...any) {
	Success(fmt.Sprintf(format, vars...))
}
