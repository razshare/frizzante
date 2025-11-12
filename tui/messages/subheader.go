package messages

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func Subheader(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	fmt.Println(configs.Styles.Subheader.Render(strings.Join(entries, "")))
}

func Subheaderf(format string, vars ...any) {
	Subheader(fmt.Sprintf(format, vars...))
}
