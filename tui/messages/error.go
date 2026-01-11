package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/tui/configs"
)

func Error(args ...any) {
	length := len(args)
	entries := make([]string, length+1)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}

	if trace := stack.Trace(); trace != "" {
		entries[length] = "\n" + stack.Trace()
	}

	Status("ERROR", strings.Join(entries, ""), configs.Colors.Error, "233", configs.Colors.Error)
}
