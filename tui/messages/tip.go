package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/config"
)

func Tip(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("TIP", strings.Join(entries, ""), config.Colors.Tip, "17", config.Colors.Tip)
}

func Tipf(format string, vars ...any) {
	Tip(fmt.Sprintf(format, vars...))
}
