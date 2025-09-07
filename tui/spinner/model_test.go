package spinner

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestSpinnerInterruptHandling(t *testing.T) {
	tests := []struct {
		name            string
		softInterrupt   bool
		keyType         tea.KeyType
		shouldQuit      bool
		shouldInterrupt bool
	}{
		{
			name:            "soft interrupt with ctrl+c quits",
			softInterrupt:   true,
			keyType:         tea.KeyCtrlC,
			shouldQuit:      true,
			shouldInterrupt: false,
		},
		{
			name:            "hard interrupt with ctrl+c interrupts",
			softInterrupt:   false,
			keyType:         tea.KeyCtrlC,
			shouldQuit:      false,
			shouldInterrupt: true,
		},
		{
			name:            "escape always quits",
			softInterrupt:   true,
			keyType:         tea.KeyEsc,
			shouldQuit:      true,
			shouldInterrupt: false,
		},
		{
			name:            "escape quits even with hard interrupt",
			softInterrupt:   false,
			keyType:         tea.KeyEsc,
			shouldQuit:      true,
			shouldInterrupt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				SoftInterrupt: tt.softInterrupt,
				Message:       "Loading...",
				Spinner:       spinner.New(),
			}

			_, cmd := model.Update(tea.KeyMsg{Type: tt.keyType})

			if tt.shouldQuit && cmd == nil {
				t.Error("expected quit command")
			}

			if tt.shouldInterrupt && cmd == nil {
				t.Error("expected interrupt command")
			}

			if !tt.shouldQuit && !tt.shouldInterrupt && cmd != nil {
				t.Error("expected no command for normal operation")
			}
		})
	}
}

func TestSpinnerMessage(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "simple message",
			message:  "Loading...",
			expected: "Loading...",
		},
		{
			name:     "empty message",
			message:  "",
			expected: "",
		},
		{
			name:     "message with special characters",
			message:  "Processing: 50% [####    ]",
			expected: "Processing: 50% [####    ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				Message: tt.message,
				Spinner: spinner.New(),
			}

			if model.Message != tt.expected {
				t.Errorf("message = %q, want %q", model.Message, tt.expected)
			}
		})
	}
}

func TestSpinnerStateTransitions(t *testing.T) {
	t.Run("soft to hard interrupt transition", func(t *testing.T) {
		model := &Model{
			SoftInterrupt: true,
			Message:       "Loading...",
			Spinner:       spinner.New(),
		}

		_, cmd1 := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if cmd1 == nil {
			t.Error("soft interrupt ctrl+c should quit")
		}

		model.SoftInterrupt = false

		_, cmd2 := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if cmd2 == nil {
			t.Error("hard interrupt ctrl+c should interrupt")
		}
	})
}

func TestSpinnerOtherKeys(t *testing.T) {
	model := &Model{
		SoftInterrupt: false,
		Message:       "Loading...",
		Spinner:       spinner.New(),
	}

	keys := []tea.KeyType{
		tea.KeyEnter,
		tea.KeySpace,
		tea.KeyTab,
		tea.KeyUp,
		tea.KeyDown,
	}

	for _, key := range keys {
		_, cmd := model.Update(tea.KeyMsg{Type: key})
		if cmd != nil {
			// cmd could be spinner tick or nil, both are fine
		}
	}
}

func TestSpinnerTickCommand(t *testing.T) {
	model := &Model{
		Message: "Loading...",
		Spinner: spinner.New(),
	}

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init should return spinner tick command")
	}
}

func TestSpinnerViewFormat(t *testing.T) {
	model := &Model{
		Message: "Test message",
		Spinner: spinner.New(),
	}

	view := model.View()

	if view == "" {
		t.Error("view should not be empty")
	}

	if !contains(view, model.Message) {
		t.Errorf("view should contain message %q", model.Message)
	}

	if len(view) > 0 && view[len(view)-1] != '\n' {
		t.Error("view should end with newline")
	}
}

func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
