package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func Section(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}
	fmt.Println(config.Styles.Section.Render("## " + strings.Join(entries, "")))
}

func Sectionf(format string, vars ...any) {
	Section(fmt.Sprintf(format, vars...))
}
