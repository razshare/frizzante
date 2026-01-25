package messages

import (
	"fmt"
	"os"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func Fatal(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	Status("ERROR", strings.Join(entries, ""), configs.Colors.Error, "233", configs.Colors.Error)
	os.Exit(1)
}
