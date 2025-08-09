package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Section(args ...any) {
	l := len(args)
	entries := make([]string, l)
	for i := 0; i < l; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	fmt.Println(config.Styles.Section.Render("## " + strings.Join(entries, "")))
}

func Sectionf(format string, vars ...any) {
	Section(fmt.Sprintf(format, vars...))
}
