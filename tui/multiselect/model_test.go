package multiselect

import (
	"slices"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestMultiselectToggleSelection(t *testing.T) {
	tests := []struct {
		name             string
		initialSelected  []string
		filteredChoices  []search.Choice
		cursorPosition   int
		expectedSelected []string
		description      string
	}{
		{
			name:            "select first item",
			initialSelected: []string{},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			cursorPosition:   0,
			expectedSelected: []string{"apple"},
			description:      "adds item when not selected",
		},
		{
			name:            "deselect already selected item",
			initialSelected: []string{"apple"},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			cursorPosition:   0,
			expectedSelected: []string{},
			description:      "removes item when already selected",
		},
		{
			name:            "select multiple items",
			initialSelected: []string{"apple"},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			cursorPosition:   1,
			expectedSelected: []string{"apple", "banana"},
			description:      "adds to existing selection",
		},
		{
			name:            "deselect from multiple",
			initialSelected: []string{"apple", "banana", "cherry"},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			cursorPosition:   1,
			expectedSelected: []string{"apple", "cherry"},
			description:      "removes middle item from selection",
		},
		{
			name:             "empty filtered list",
			initialSelected:  []string{},
			filteredChoices:  []search.Choice{},
			cursorPosition:   0,
			expectedSelected: []string{},
			description:      "no change when filtered list is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				Selected: tt.initialSelected,
				Search: &search.Search{
					Filtered: tt.filteredChoices,
					Input:    textinput.New(),
				},
				Viewport: &viewport.Viewport{
					Cursor: tt.cursorPosition,
				},
			}

			model.Update(tea.KeyMsg{Type: tea.KeySpace})

			if !slices.Equal(model.Selected, tt.expectedSelected) {
				t.Errorf("%s: got %v, want %v", tt.description, model.Selected, tt.expectedSelected)
			}
		})
	}
}

func TestMultiselectEnterBehavior(t *testing.T) {
	tests := []struct {
		name             string
		initialSelected  []string
		filteredChoices  []search.Choice
		cursorPosition   int
		expectedSelected []string
		shouldQuit       bool
	}{
		{
			name:            "enter with existing selection",
			initialSelected: []string{"apple", "banana"},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			cursorPosition:   2,
			expectedSelected: []string{"apple", "banana"},
			shouldQuit:       true,
		},
		{
			name:            "enter with no selection auto-selects current",
			initialSelected: []string{},
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			cursorPosition:   1,
			expectedSelected: []string{"banana"},
			shouldQuit:       true,
		},
		{
			name:             "enter with empty filtered list",
			initialSelected:  []string{},
			filteredChoices:  []search.Choice{},
			cursorPosition:   0,
			expectedSelected: []string{},
			shouldQuit:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				Selected: tt.initialSelected,
				Search: &search.Search{
					Filtered: tt.filteredChoices,
					Input:    textinput.New(),
				},
				Viewport: &viewport.Viewport{
					Cursor: tt.cursorPosition,
				},
			}

			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			if !slices.Equal(model.Selected, tt.expectedSelected) {
				t.Errorf("got selected %v, want %v", model.Selected, tt.expectedSelected)
			}

			if tt.shouldQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			}
		})
	}
}

func TestMultiselectEscapeBehavior(t *testing.T) {
	t.Run("escape clears search when active", func(t *testing.T) {
		model := &Model{
			Selected: []string{"apple"},
			Search: &search.Search{
				Active: true,
				Choices: []search.Choice{
					{Id: "apple"},
					{Id: "banana"},
				},
				Filtered: []search.Choice{
					{Id: "apple"},
				},
				Input: textinput.New(),
			},
			Viewport: &viewport.Viewport{
				Cursor: 0,
			},
		}

		model.Search.Input.SetValue("app")

		model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if model.Search.Active {
			t.Error("search should be inactive after escape")
		}

		if model.Search.Input.Value() != "" {
			t.Error("search input should be cleared")
		}

		if len(model.Search.Filtered) != len(model.Search.Choices) {
			t.Error("filtered should be reset to all choices")
		}
	})

	t.Run("escape quits when search inactive", func(t *testing.T) {
		model := &Model{
			Selected: []string{"apple"},
			Search: &search.Search{
				Active: false,
				Input:  textinput.New(),
			},
			Viewport: &viewport.Viewport{},
		}

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if len(model.Selected) != 0 {
			t.Error("selected should be cleared on quit")
		}

		if cmd == nil {
			t.Error("expected quit command")
		}
	})
}

func TestMultiselectNavigation(t *testing.T) {
	choices := []search.Choice{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
	}

	tests := []struct {
		name           string
		keyType        tea.KeyType
		initialCursor  int
		expectedCursor int
	}{
		{"move down with arrow", tea.KeyDown, 0, 1},
		{"move down with tab", tea.KeyTab, 0, 1},
		{"move down with ctrl+n", tea.KeyCtrlN, 0, 1},
		{"move up with arrow", tea.KeyUp, 2, 1},
		{"move up with ctrl+p", tea.KeyCtrlP, 2, 1},
		{"wrap from bottom", tea.KeyDown, 4, 0},
		{"wrap from top", tea.KeyUp, 0, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				Search: &search.Search{
					Filtered: choices,
					Input:    textinput.New(),
				},
				Viewport: &viewport.Viewport{
					Cursor:  tt.initialCursor,
					Visible: 5,
				},
			}

			model.Update(tea.KeyMsg{Type: tt.keyType})

			if model.Viewport.Cursor != tt.expectedCursor {
				t.Errorf("cursor = %d, want %d", model.Viewport.Cursor, tt.expectedCursor)
			}
		})
	}
}

func TestMultiselectSearchActivation(t *testing.T) {
	model := &Model{
		Search: &search.Search{
			Active: false,
			Choices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			Filtered: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{},
	}

	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	if !model.Search.Active {
		t.Error("search should be activated on character input")
	}

	if model.Search.Input.Value() != "a" {
		t.Errorf("search input = %q, want 'a'", model.Search.Input.Value())
	}
}

func TestMultiselectSelectionPersistence(t *testing.T) {
	model := &Model{
		Selected: []string{},
		Search: &search.Search{
			Filtered: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{
			Cursor:  0,
			Visible: 3,
		},
	}

	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	model.Viewport.Cursor = 2
	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	expected := []string{"apple", "cherry"}
	if !slices.Equal(model.Selected, expected) {
		t.Errorf("selected = %v, want %v", model.Selected, expected)
	}

	model.Viewport.Cursor = 0
	model.Update(tea.KeyMsg{Type: tea.KeySpace})

	expected = []string{"cherry"}
	if !slices.Equal(model.Selected, expected) {
		t.Errorf("after deselect, selected = %v, want %v", model.Selected, expected)
	}
}
