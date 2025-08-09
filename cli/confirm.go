package cli

import (
	"fmt"
	"github.com/razshare/frizzante/tui/confirm"
)

// Confirm shows a confirmation prompt.
//
// Returns true if the user confirms, otherwise false.
func Confirm(text string) bool {
	if *FlagYes {
		return true
	}

	yes, showError := confirm.Send(text, true)

	if showError != nil {
		Fatal(showError)
	}

	return yes
}

// Confirmf shows a confirmation prompt.
//
// Returns true if the user confirms, otherwise false.
func Confirmf(template string, vars ...any) bool {
	return Confirm(fmt.Sprintf(template, vars...))
}
