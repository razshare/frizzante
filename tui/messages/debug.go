package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func Debug(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("DEBUG", strings.Join(entries, ""), configs.Colors.Debug, "17", configs.Colors.Debug)
}
