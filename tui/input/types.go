package input

import "github.com/charmbracelet/bubbles/textinput"

type Model struct {
	TextInput textinput.Model
	Prompt    string
}
