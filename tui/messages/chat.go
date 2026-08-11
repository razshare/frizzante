package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func Chat(user string, args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status(user, strings.Join(entries, ""), configs.Colors.Input, "17", configs.Colors.Input)
}
