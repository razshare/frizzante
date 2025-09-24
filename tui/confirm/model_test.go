package confirm

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate(t *testing.T) {
	var model *Model
	var cmd tea.Cmd

	// y key confirms
	model = &Model{Prompt: "Test prompt", DefaultValue: false, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if !model.Confirmed {
		t.Fatal("confirm should be true")
	}

	// Y key confirms (case-insensitive)
	model = &Model{Prompt: "Test prompt", DefaultValue: false, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
	if !model.Confirmed {
		t.Fatal("confirm should be true")
	}

	// n key confirms
	model = &Model{Prompt: "Test prompt", DefaultValue: false, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if model.Confirmed {
		t.Fatal("confirm should not be true")
	}

	// N key confirms (case-insensitive)
	model = &Model{Prompt: "Test prompt", DefaultValue: false, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	if model.Confirmed {
		t.Fatal("confirm should not be true")
	}

	// enter uses default true
	model = &Model{Prompt: "Test prompt", DefaultValue: true, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !model.Confirmed {
		t.Fatal("confirm should be true")
	}

	// enter uses default false
	model = &Model{Prompt: "Test prompt", DefaultValue: false, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if model.Confirmed {
		t.Fatal("confirm should not be true")
	}

	// escape denies
	model = &Model{Prompt: "Test prompt", DefaultValue: true, Confirmed: false}
	model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.Confirmed {
		t.Fatal("confirm should not be true")
	}

	// ctrl+c interrupts
	model = &Model{Prompt: "Test prompt", DefaultValue: true, Confirmed: false}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if model.Confirmed {
		t.Fatal("confirm should not be true")
	}
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("confirm should interrupt")
	}

	// other key does nothing
	model = &Model{Prompt: "Test prompt", DefaultValue: true, Confirmed: false}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	if cmd != nil {
		t.Fatal("confirm should ignore irrelevant keys")
	}
}
