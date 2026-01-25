package inputs

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestInputEscapeBehavior(t *testing.T) {
	t.Run("escape clears input when has value", func(t *testing.T) {
		model := &Model{
			Prompt: "Enter name:",
		}
		model.Value = "test value"
		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if model.Value != "" {
			t.Errorf("input value = %q, want empty", model.Value)
		}
		if cmd != nil {
			t.Error("should not quit when clearing input")
		}
	})
	t.Run("escape quits when input empty", func(t *testing.T) {
		model := &Model{
			Prompt: "Enter name:",
		}
		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		if cmd == nil {
			t.Error("should quit when input is empty")
		}
	})
	t.Run("escape clears then quits on second press", func(t *testing.T) {
		model := &Model{
			Prompt: "Enter name:",
		}
		model.Value = "test"
		_, cmd1 := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		if cmd1 != nil {
			t.Error("first escape should not quit")
		}
		if model.Value != "" {
			t.Error("first escape should clear input")
		}
		_, cmd2 := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		if cmd2 == nil {
			t.Error("second escape should quit")
		}
	})
}

func TestInputEnterBehavior(t *testing.T) {
	type Input struct {
		name  string
		value string
	}
	input := Input{"enter with value", "test input"}
	model := &Model{Prompt: "Enter value:"}
	model.Value = input.value
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("enter should always quit")
	}
	if model.Value != input.value {
		t.Errorf("value should be preserved, got %q, want %q", model.Value, input.value)
	}
	input = Input{"enter with empty", ""}
	model = &Model{Prompt: "Enter value:"}
	model.Value = input.value
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("enter should always quit")
	}
	if model.Value != input.value {
		t.Errorf("value should be preserved, got %q, want %q", model.Value, input.value)
	}
	input = Input{"enter with spaces", "   "}
	model = &Model{Prompt: "Enter value:"}
	model.Value = input.value
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("enter should always quit")
	}
	if model.Value != input.value {
		t.Errorf("value should be preserved, got %q, want %q", model.Value, input.value)
	}
}

func TestInputInterrupt(t *testing.T) {
	model := &Model{
		Prompt: "Enter value:",
	}
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("ctrl+c should send interrupt")
	}
}

func TestInputTextEntry(t *testing.T) {
	type Input struct {
		key      tea.KeyMsg
		expected string
	}
	model := &Model{Prompt: "Enter name:"}
	inputs := []Input{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, "h"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}, "he"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, "hel"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, "hell"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}}, "hello"},
	}
	for _, input := range inputs {
		model.Update(input.key)
		if model.Value != input.expected {
			t.Errorf("after input, value = %q, want %q", model.Value, input.expected)
		}
	}
}

func TestInputBackspace(t *testing.T) {
	model := &Model{Prompt: "Enter value:"}
	model.Value = "hello"
	model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if model.Value != "hell" {
		t.Errorf("after backspace, value = %q, want 'hell'", model.Value)
	}
}

func TestInputStatePreservation(t *testing.T) {
	model := &Model{Prompt: "Enter your name:"}
	originalPrompt := model.Prompt
	model.Value = "test"
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if model.Prompt != originalPrompt {
		t.Error("prompt should not change during input")
	}
}
