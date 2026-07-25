package select_one

import (
	"fmt"

	"github.com/razshare/frizzante/v2/tui/search"
)

func Sendf(choices []search.Choice, format string, vars ...any) (selected string, err error) {
	return Send(choices, fmt.Sprintf(format, vars...))
}
