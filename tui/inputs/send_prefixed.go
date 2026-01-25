package inputs

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
)

func SendPrefixed(message string, prefix string) (value string, err error) {
	input := textinput.New()
	input.Placeholder = "Type here..."
	input.Focus()
	input.Width = 50
	var model *Model
	if model, err = program.Run(&Model{Prompt: message, Prefix: prefix}); err != nil {
		return
	}
	value = model.Prefix + model.Value
	return
}
