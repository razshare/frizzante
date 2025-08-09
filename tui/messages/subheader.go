package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Subheader(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	fmt.Println(config.Styles.Subheader.Render(strings.Join(entries, "")))
}

func Subheaderf(format string, vars ...any) {
	Subheader(fmt.Sprintf(format, vars...))
}
