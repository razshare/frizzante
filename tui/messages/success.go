package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Success(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("SUCCESS", strings.Join(entries, ""), config.Colors.Success, "17", config.Colors.Success)
}

func Successf(format string, vars ...any) {
	Success(fmt.Sprintf(format, vars...))
}
