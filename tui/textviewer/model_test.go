package textviewer

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestUpdate(t *testing.T) {
	var model *Model
	var cmd tea.Cmd

	// ctrl+c
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should interrupt")
	}

	// escape clears search when active
	model = &Model{
		Search: &search.Search{
			Active:   true,
			Choices:  []search.Choice{{Id: "apple"}, {Id: "banana"}},
			Filtered: []search.Choice{{Id: "apple"}},
			Input:    textinput.New(),
		},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Search.Input.SetValue("app")
	model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.Search.Active {
		t.Fatal("singleselect search be inactive")
	}
	if model.Search.Input.Value() != "" {
		t.Fatal("singleselect search should be empty")
	}
	if len(model.Search.Filtered) != len(model.Search.Choices) {
		t.Fatal("singleselect choices should be the same as filtered choices")
	}

	// escape quits when search inactive
	model = &Model{
		Search:   &search.Search{Active: false, Input: textinput.New()},
		Viewport: &viewport.Viewport{},
	}

	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should quit")
	}

	// move down with arrow
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with tab
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyTab})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with ctrl+pgdown
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyCtrlPgDown})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with arrow
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with tab
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with ctrl+pgdown
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyCtrlPgUp})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down and wrap
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 4},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.Viewport.Cursor != 0 {
		t.Fatal("multi select cursor should be 0")
	}

	// move up and wrap
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}, Input: textinput.New()},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.Viewport.Cursor != 4 {
		t.Fatal("multi select cursor should be 4")
	}

	// activate search when typing
	model = &Model{
		Search: &search.Search{
			Active:   false,
			Choices:  []search.Choice{{Id: "apple"}, {Id: "banana"}},
			Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}},
			Input:    textinput.New(),
		},
		Viewport: &viewport.Viewport{},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if !model.Search.Active {
		t.Fatal("singleselect search should be active")
	}
	if model.Search.Input.Value() != "a" {
		t.Fatal("singleselect search value should be a")
	}
}
