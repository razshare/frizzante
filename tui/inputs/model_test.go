package inputs

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestInputEscapeBehavior(t *testing.T) {
	t.Run("escape clears input when has value", func(t *testing.T) {
		model := &Model{
			Prompt:    "Enter name:",
			TextInput: textinput.New(),
		}
		model.TextInput.SetValue("test value")

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if model.TextInput.Value() != "" {
			t.Errorf("input value = %q, want empty", model.TextInput.Value())
		}

		if cmd != nil {
			t.Error("should not quit when clearing input")
		}
	})

	t.Run("escape quits when input empty", func(t *testing.T) {
		model := &Model{
			Prompt:    "Enter name:",
			TextInput: textinput.New(),
		}

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if cmd == nil {
			t.Error("should quit when input is empty")
		}
	})

	t.Run("escape clears then quits on second press", func(t *testing.T) {
		model := &Model{
			Prompt:    "Enter name:",
			TextInput: textinput.New(),
		}
		model.TextInput.SetValue("test")

		_, cmd1 := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		if cmd1 != nil {
			t.Error("first escape should not quit")
		}

		if model.TextInput.Value() != "" {
			t.Error("first escape should clear input")
		}

		_, cmd2 := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		if cmd2 == nil {
			t.Error("second escape should quit")
		}
	})
}

func TestInputEnterBehavior(t *testing.T) {
	type TestData struct {
		name       string
		inputValue string
	}

	data := []TestData{
		{"enter with value", "test input"},
		{"enter with empty", ""},
		{"enter with spaces", "   "},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			model := &Model{
				Prompt:    "Enter value:",
				TextInput: textinput.New(),
			}
			model.TextInput.SetValue(d.inputValue)

			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			if cmd == nil {
				t.Error("enter should always quit")
			}

			if model.TextInput.Value() != d.inputValue {
				t.Errorf("value should be preserved, got %q, want %q", model.TextInput.Value(), d.inputValue)
			}
		})
	}
}

func TestInputInterrupt(t *testing.T) {
	model := &Model{
		Prompt:    "Enter value:",
		TextInput: textinput.New(),
	}

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("ctrl+c should send interrupt")
	}
}

func TestInputTextEntry(t *testing.T) {
	model := &Model{
		Prompt:    "Enter name:",
		TextInput: textinput.New(),
	}
	model.TextInput.Focus()

	inputs := []struct {
		key      tea.KeyMsg
		expected string
	}{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, "h"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}, "he"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, "hel"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, "hell"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}, "hello"},
	}

	for _, input := range inputs {
		model.Update(input.key)
		if model.TextInput.Value() != input.expected {
			t.Errorf("after input, value = %q, want %q", model.TextInput.Value(), input.expected)
		}
	}
}

func TestInputBackspace(t *testing.T) {
	model := &Model{
		Prompt:    "Enter value:",
		TextInput: textinput.New(),
	}
	model.TextInput.SetValue("hello")
	model.TextInput.Focus()

	model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	if model.TextInput.Value() != "hell" {
		t.Errorf("after backspace, value = %q, want 'hell'", model.TextInput.Value())
	}
}

func TestInputStatePreservation(t *testing.T) {
	model := &Model{
		Prompt:    "Enter your name:",
		TextInput: textinput.New(),
	}
	originalPrompt := model.Prompt

	model.TextInput.SetValue("test")
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	if model.Prompt != originalPrompt {
		t.Error("prompt should not change during input")
	}
}

func TestInputResetFunctionality(t *testing.T) {
	model := &Model{
		Prompt:    "Enter value:",
		TextInput: textinput.New(),
	}

	model.TextInput.SetValue("some text")
	model.TextInput.SetCursor(5)

	model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if model.TextInput.Value() != "" {
		t.Error("reset should clear value")
	}

	if model.TextInput.Position() != 0 {
		t.Error("reset should reset cursor position")
	}
}
