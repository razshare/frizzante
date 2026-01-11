package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func Success(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("SUCCESS", strings.Join(entries, ""), configs.Colors.Success, "17", configs.Colors.Success)
}
