package action

import (
	"github.com/razshare/frizzante/internal/tui/text"
)

func Clear(_ ClearOptions) error {
	text.Clrscr()
	return nil
}
