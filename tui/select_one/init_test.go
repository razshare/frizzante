package select_one

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestUpdate(t *testing.T) {
	var model *Model
	var cmd tea.Cmd

	// select first item
	model = &Model{
		Selected: "",
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if model.Selected != "apple" {
		t.Fatal("singleselect should be apple")
	}
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should quit")
	}

	// move down then select
	model = &Model{
		Selected: "",
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if model.Selected != "banana" {
		t.Fatal("singleselect should be banana")
	}
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should quit")
	}

	// ctrl+c
	model = &Model{
		Selected: "",
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if model.Selected != "" {
		t.Fatal("singleselect should be empty")
	}
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should interrupt")
	}

	// escape clears search when active
	model = &Model{
		Selected: "",
		Search: &search.Search{
			Active:   true,
			Choices:  []search.Choice{{Id: "apple"}, {Id: "banana"}},
			Filtered: []search.Choice{{Id: "apple"}},
		},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Search.Value = "app"
	model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if model.Search.Active {
		t.Fatal("singleselect search be inactive")
	}
	if model.Search.Value != "" {
		t.Fatal("singleselect search should be empty")
	}
	if len(model.Search.Filtered) != len(model.Search.Choices) {
		t.Fatal("singleselect choices should be the same as filtered choices")
	}

	// escape quits when search inactive
	model = &Model{
		Selected: "apple",
		Search:   &search.Search{Active: false},
		Viewport: &viewport.Viewport{},
	}

	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if len(model.Selected) != 0 {
		t.Fatal("singleselect should be empty")
	}
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("singleselect should quit")
	}

	// move down with arrow
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with tab
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyTab})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with ctrl+pgdown
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyCtrlPgDown})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with arrow
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyUp})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with tab
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down with ctrl+pgdown
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 2},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyCtrlPgUp})
	if model.Viewport.Cursor != 1 {
		t.Fatal("multi select cursor should be 1")
	}

	// move down and wrap
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
		Viewport: &viewport.Viewport{Visible: 5, Cursor: 4},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	if model.Viewport.Cursor != 0 {
		t.Fatal("multi select cursor should be 0")
	}

	// move up and wrap
	model = &Model{
		Search:   &search.Search{Filtered: []search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}},
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
		},
		Viewport: &viewport.Viewport{},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if !model.Search.Active {
		t.Fatal("singleselect search should be active")
	}
	if model.Search.Value != "a" {
		t.Fatal("singleselect search value should be a")
	}

	// moving to beginning of viewport with home
	model = &Model{
		Search: &search.Search{
			Active:   false,
			Choices:  []search.Choice{{Id: "apple"}, {Id: "banana1"}, {Id: "banana2"}, {Id: "banana3"}, {Id: "banana4"}, {Id: "banana5"}},
			Filtered: []search.Choice{{Id: "apple"}, {Id: "banana1"}, {Id: "banana2"}, {Id: "banana3"}, {Id: "banana4"}, {Id: "banana5"}},
		},
		Viewport: &viewport.Viewport{
			Offset:  1,
			Cursor:  6,
			Visible: 5,
		},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyEnd})
	if model.Viewport.Cursor != 5 {
		t.Fatal("multi select cursor should be 6")
	}

	// moving to end of viewport with end
	model = &Model{
		Search: &search.Search{
			Active:   false,
			Choices:  []search.Choice{{Id: "apple"}, {Id: "banana1"}, {Id: "banana2"}, {Id: "banana3"}, {Id: "banana4"}, {Id: "banana5"}},
			Filtered: []search.Choice{{Id: "apple"}, {Id: "banana1"}, {Id: "banana2"}, {Id: "banana3"}, {Id: "banana4"}, {Id: "banana5"}},
		},
		Viewport: &viewport.Viewport{
			Offset:  0,
			Cursor:  5,
			Visible: 5,
		},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyHome})
	if model.Viewport.Cursor != 0 {
		t.Fatal("multi select cursor should be 0")
	}
}
