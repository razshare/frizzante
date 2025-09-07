package singleselect

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestSingleSelectSelection(t *testing.T) {
	tests := []struct {
		name             string
		filteredChoices  []search.Choice
		cursorPosition   int
		expectedSelected string
		shouldQuit       bool
	}{
		{
			name: "select first item",
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			cursorPosition:   0,
			expectedSelected: "apple",
			shouldQuit:       true,
		},
		{
			name: "select middle item",
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			cursorPosition:   1,
			expectedSelected: "banana",
			shouldQuit:       true,
		},
		{
			name: "select last item",
			filteredChoices: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
				{Id: "cherry"},
			},
			cursorPosition:   2,
			expectedSelected: "cherry",
			shouldQuit:       true,
		},
		{
			name:             "empty filtered list",
			filteredChoices:  []search.Choice{},
			cursorPosition:   0,
			expectedSelected: "",
			shouldQuit:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{
				Selected: "",
				Search: &search.Search{
					Filtered: tt.filteredChoices,
					Input:    textinput.New(),
				},
				Viewport: &viewport.Viewport{
					Cursor: tt.cursorPosition,
				},
			}

			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			if model.Selected != tt.expectedSelected {
				t.Errorf("got selected %q, want %q", model.Selected, tt.expectedSelected)
			}

			if tt.shouldQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			} else if !tt.shouldQuit && cmd != nil {
				t.Error("expected no command, got quit")
			}
		})
	}
}

func TestSingleSelectEscapeBehavior(t *testing.T) {
	t.Run("escape clears search when active", func(t *testing.T) {
		model := &Model{
			Selected: "",
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

	t.Run("escape quits and clears selection when search inactive", func(t *testing.T) {
		model := &Model{
			Selected: "previous",
			Search: &search.Search{
				Active: false,
				Input:  textinput.New(),
			},
			Viewport: &viewport.Viewport{},
		}

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if model.Selected != "" {
			t.Errorf("selected should be cleared, got %q", model.Selected)
		}

		if cmd == nil {
			t.Error("expected quit command")
		}
	})
}

func TestSingleSelectNavigation(t *testing.T) {
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

func TestSingleSelectReplaceSelection(t *testing.T) {
	model := &Model{
		Selected: "initial",
		Search: &search.Search{
			Filtered: []search.Choice{
				{Id: "apple"},
				{Id: "banana"},
			},
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{
			Cursor: 1,
		},
	}

	model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if model.Selected != "banana" {
		t.Errorf("selected = %q, want 'banana'", model.Selected)
	}
}

func TestSingleSelectInterrupt(t *testing.T) {
	model := &Model{
		Search: &search.Search{
			Input: textinput.New(),
		},
		Viewport: &viewport.Viewport{},
	}

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Error("expected interrupt command")
	}
}
