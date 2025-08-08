package cli

import (
	"fmt"
	"github.com/pterm/pterm"
)

// Confirm shows a confirmation prompt.
//
// Returns true if the user confirms, otherwise false.
func Confirm(text string) bool {
	if *FlagYes {
		return true
	}

	yes, showError := pterm.
		DefaultInteractiveConfirm.
		WithConfirmText("Y").
		WithDefaultText("n").
		WithDefaultValue(true).
		Show(text)

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
