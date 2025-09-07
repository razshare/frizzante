package npmselect

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestNPMSelectDebouncing(t *testing.T) {
	model := &Model{
		Selected: []string{},
		Search: &search.Search{
			Active: false,
			Input:  textinput.New(),
		},
		Viewport:  &viewport.Viewport{},
		Debouncer: time.NewTimer(100 * time.Millisecond),
		Debounce:  100 * time.Millisecond,
	}

	model.Debouncer.Stop()

	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	if !model.Search.Active {
		t.Error("search should be activated")
	}

	if model.Search.Input.Value() != "r" {
		t.Errorf("input value = %q, want 'r'", model.Search.Input.Value())
	}

	if model.LastQuery != "r" {
		t.Errorf("last query = %q, want 'r'", model.LastQuery)
	}
}

func TestNPMSelectPackageTransformation(t *testing.T) {
	packages := []npm.PackageInfo{
		{
			Name:        "react",
			Version:     "18.2.0",
			Description: "A JavaScript library for building user interfaces",
		},
		{
			Name:        "vue",
			Version:     "3.3.0",
			Description: "This is a very long description that should be truncated because it exceeds the maximum allowed length for display",
		},
		{
			Name:        "angular",
			Version:     "16.0.0",
			Description: "",
		},
	}

	choices := make([]search.Choice, 0, len(packages))
	for _, pkg := range packages {
		id := fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
		desc := pkg.Description
		if len(desc) > 50 {
			desc = desc[:50]
		}
		choices = append(choices, search.Choice{
			Id:          id,
			Description: desc,
		})
	}

	expected := []struct {
		id   string
		desc string
	}{
		{"react@18.2.0", "A JavaScript library for building user interfaces"},
		{"vue@3.3.0", "This is a very long description that should be tru"},
		{"angular@16.0.0", ""},
	}

	for i, choice := range choices {
		if choice.Id != expected[i].id {
			t.Errorf("choice[%d].Id = %q, want %q", i, choice.Id, expected[i].id)
		}
		if choice.Description != expected[i].desc {
			t.Errorf("choice[%d].Description = %q, want %q", i, choice.Description, expected[i].desc)
		}
	}
}

func TestNPMSelectSelection(t *testing.T) {
	model := &Model{
		Selected: []string{},
		Search: &search.Search{
			Filtered: []search.Choice{
				{Id: "react@18.2.0"},
				{Id: "vue@3.3.0"},
				{Id: "angular@16.0.0"},
			},
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{
			Cursor: 0,
		},
	}

	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	if !slices.Contains(model.Selected, "react@18.2.0") {
		t.Error("react@18.2.0 should be selected")
	}

	model.Viewport.Cursor = 1
	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	if len(model.Selected) != 2 {
		t.Errorf("selected count = %d, want 2", len(model.Selected))
	}

	model.Viewport.Cursor = 0
	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	if slices.Contains(model.Selected, "react@18.2.0") {
		t.Error("react@18.2.0 should be deselected")
	}

	if len(model.Selected) != 1 {
		t.Errorf("selected count = %d, want 1", len(model.Selected))
	}
}

func TestNPMSelectAutoSelection(t *testing.T) {
	model := &Model{
		Selected: []string{},
		Search: &search.Search{
			Filtered: []search.Choice{
				{Id: "react@18.2.0"},
				{Id: "vue@3.3.0"},
			},
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{
			Cursor: 1,
		},
	}

	model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !model.Confirmed {
		t.Error("should be confirmed")
	}

	if len(model.Selected) != 1 {
		t.Errorf("selected count = %d, want 1", len(model.Selected))
	}

	if model.Selected[0] != "vue@3.3.0" {
		t.Errorf("selected = %q, want 'vue@3.3.0'", model.Selected[0])
	}
}

func TestNPMSelectQuitBehavior(t *testing.T) {
	t.Run("escape quits and clears", func(t *testing.T) {
		model := &Model{
			Selected: []string{"react@18.2.0"},
			Search: &search.Search{
				Input: textinput.New(),
			},
			Viewport: &viewport.Viewport{},
		}

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if !model.Quitting {
			t.Error("should be quitting")
		}

		if len(model.Selected) != 0 {
			t.Error("selected should be cleared")
		}

		if cmd == nil {
			t.Error("expected quit command")
		}
	})

	t.Run("enter confirms and quits", func(t *testing.T) {
		model := &Model{
			Selected: []string{"react@18.2.0"},
			Search: &search.Search{
				Filtered: []search.Choice{},
				Input:    textinput.New(),
			},
			Viewport: &viewport.Viewport{},
		}

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		if !model.Confirmed {
			t.Error("should be confirmed")
		}

		if cmd == nil {
			t.Error("expected quit command")
		}
	})
}

func TestNPMSelectLoadingStates(t *testing.T) {
	model := &Model{
		Loading: false,
		Error:   nil,
		Search: &search.Search{
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{},
	}

	model.Loading = true
	if !model.Loading {
		t.Error("should be in loading state")
	}

	model.Loading = false
	model.Error = fmt.Errorf("network error")

	if model.Error == nil {
		t.Error("should have error state")
	}

	if model.Error.Error() != "network error" {
		t.Errorf("error = %q, want 'network error'", model.Error.Error())
	}
}

func TestNPMSelectSearchResultProcessing(t *testing.T) {
	msg := SearchResultMsg{
		Packages: []npm.PackageInfo{
			{Name: "react", Version: "18.2.0", Description: "React library"},
			{Name: "react-dom", Version: "18.2.0", Description: "React DOM"},
		},
		Error: nil,
	}

	if msg.Error != nil {
		t.Error("should not have error")
	}

	if len(msg.Packages) != 2 {
		t.Errorf("packages count = %d, want 2", len(msg.Packages))
	}

	errorMsg := SearchResultMsg{
		Packages: []npm.PackageInfo{},
		Error:    fmt.Errorf("timeout"),
	}

	if errorMsg.Error == nil {
		t.Error("should have error")
	}

	if len(errorMsg.Packages) != 0 {
		t.Error("should have empty packages on error")
	}
}
