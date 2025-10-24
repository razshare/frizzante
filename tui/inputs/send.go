package inputs

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
)

func Send(message string) (value string, err error) {
	input := textinput.New()
	input.Placeholder = "Type here..."
	input.Focus()
	input.Width = 50

	var model *Model
	if model, err = program.Run(&Model{TextInput: input, Prompt: message}); err != nil {
		return
	}

	value = model.TextInput.Value()

	return
}

func Sendf(format string, vars ...any) (value string, err error) {
	return Send(fmt.Sprintf(format, vars...))
}
