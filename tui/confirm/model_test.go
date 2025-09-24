package confirm

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConfirmKeyHandling(t *testing.T) {
	type TestData struct {
		name            string
		defaultValue    bool
		keyInput        tea.KeyMsg
		expectedConfirm bool
		shouldQuit      bool
		shouldInterrupt bool
	}

	data := []TestData{
		{
			name:            "y key confirms",
			defaultValue:    false,
			keyInput:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			expectedConfirm: true,
			shouldQuit:      true,
		},
		{
			name:            "Y key confirms (case insensitive)",
			defaultValue:    false,
			keyInput:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}},
			expectedConfirm: true,
			shouldQuit:      true,
		},
		{
			name:            "n key denies",
			defaultValue:    true,
			keyInput:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			expectedConfirm: false,
			shouldQuit:      true,
		},
		{
			name:            "N key denies (case insensitive)",
			defaultValue:    true,
			keyInput:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}},
			expectedConfirm: false,
			shouldQuit:      true,
		},
		{
			name:            "enter uses default true",
			defaultValue:    true,
			keyInput:        tea.KeyMsg{Type: tea.KeyEnter},
			expectedConfirm: true,
			shouldQuit:      true,
		},
		{
			name:            "enter uses default false",
			defaultValue:    false,
			keyInput:        tea.KeyMsg{Type: tea.KeyEnter},
			expectedConfirm: false,
			shouldQuit:      true,
		},
		{
			name:            "escape denies",
			defaultValue:    true,
			keyInput:        tea.KeyMsg{Type: tea.KeyEsc},
			expectedConfirm: false,
			shouldQuit:      true,
		},
		{
			name:            "ctrl+c interrupts",
			defaultValue:    false,
			keyInput:        tea.KeyMsg{Type: tea.KeyCtrlC},
			expectedConfirm: false,
			shouldInterrupt: true,
		},
		{
			name:            "other key does nothing",
			defaultValue:    false,
			keyInput:        tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
			expectedConfirm: false,
			shouldQuit:      false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			model := &Model{
				Prompt:       "Test prompt",
				DefaultValue: d.defaultValue,
				Confirmed:    false,
			}

			_, cmd := model.Update(d.keyInput)

			if model.Confirmed != d.expectedConfirm {
				t.Errorf("Confirmed = %v, want %v", model.Confirmed, d.expectedConfirm)
			}

			if d.shouldQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			}

			if d.shouldInterrupt && cmd == nil {
				t.Error("expected interrupt command, got nil")
			}

			if !d.shouldQuit && !d.shouldInterrupt && cmd != nil {
				t.Error("expected no command, got one")
			}
		})
	}
}

func TestConfirmDefaultValueBehavior(t *testing.T) {
	t.Run("default true with enter", func(t *testing.T) {
		model := &Model{
			DefaultValue: true,
			Confirmed:    false,
		}

		model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		if !model.Confirmed {
			t.Error("should be confirmed when default is true and enter pressed")
		}
	})

	t.Run("default false with enter", func(t *testing.T) {
		model := &Model{
			DefaultValue: false,
			Confirmed:    true,
		}

		model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		if model.Confirmed {
			t.Error("should not be confirmed when default is false and enter pressed")
		}
	})
}

func TestConfirmEscapeAlwaysDenies(t *testing.T) {
	type TestData struct {
		name         string
		defaultValue bool
		initial      bool
	}

	data := []TestData{
		{"escape with default true", true, true},
		{"escape with default false", false, true},
		{"escape already false", false, false},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			model := &Model{
				DefaultValue: d.defaultValue,
				Confirmed:    d.initial,
			}

			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			if model.Confirmed {
				t.Error("escape should always set Confirmed to false")
			}
		})
	}
}

func TestConfirmCaseInsensitivity(t *testing.T) {
	inputs := []struct {
		input    rune
		expected bool
	}{
		{'y', true},
		{'Y', true},
		{'n', false},
		{'N', false},
	}

	for _, tt := range inputs {
		t.Run(string(tt.input), func(t *testing.T) {
			model := &Model{}

			model.Update(tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{tt.input},
			})

			if model.Confirmed != tt.expected {
				t.Errorf("input %c: Confirmed = %v, want %v", tt.input, model.Confirmed, tt.expected)
			}
		})
	}
}

func TestConfirmStatePreservation(t *testing.T) {
	model := &Model{
		Prompt:       "Are you sure?",
		DefaultValue: true,
		Confirmed:    false,
	}

	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	if model.Prompt != "Are you sure?" {
		t.Error("prompt should not change")
	}

	if model.DefaultValue != true {
		t.Error("default value should not change")
	}

	if model.Confirmed != false {
		t.Error("confirmed should not change for invalid input")
	}
}
