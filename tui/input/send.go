package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"main/program"
)

func Send(prompt string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()
	ti.Width = 50
	m := Model{TextInput: ti, Prompt: prompt}
	result, err := program.Run(m)
	if err != nil {
		return "", err
	}
	if result.Cancelled {
		return "", fmt.Errorf("input cancelled")
	}
	return result.TextInput.Value(), nil
}
