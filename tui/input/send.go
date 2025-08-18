package input

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/tui/program"
)

func Send(prompt string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50
	model, err := program.Run(&Model{TextInput: ti, Prompt: prompt})
	if err != nil {
		return "", errors.New(err.Error() + "\n" + stack.Trace())
	}

	return model.TextInput.Value(), nil
}

func Sendf(format string, vars ...any) (string, error) {
	return Send(fmt.Sprintf(format, vars...))
}
