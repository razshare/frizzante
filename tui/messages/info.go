package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func Info(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("INFO", strings.Join(entries, ""), configs.Colors.Info, "17", configs.Colors.Info)
}
