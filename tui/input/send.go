package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
)

func Send(prompt string) string {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50
	result, err := program.Run(&Model{TextInput: ti, Prompt: prompt})
	if err != nil {
		messages.Fatal(err)
	}

	return result.TextInput.Value()
}

func Sendf(format string, vars ...any) string {
	return Send(fmt.Sprintf(format, vars...))
}
