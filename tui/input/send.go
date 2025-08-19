package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
)

func Send(msg string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50

	model, err := program.Run(&Model{TextInput: ti, Prompt: msg})

	if err != nil {
		return "", err
	}

	return model.TextInput.Value(), nil
}

func Sendf(format string, vars ...any) (string, error) {
	return Send(fmt.Sprintf(format, vars...))
}
