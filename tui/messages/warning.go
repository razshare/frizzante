package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func Warning(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("WARNING", strings.Join(entries, ""), configs.Colors.Warning, "17", configs.Colors.Warning)
}
