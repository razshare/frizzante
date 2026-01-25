package select_many

import (
	"reflect"
	"slices"
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
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{"apple"}) {
		t.Fatal("multiselect should contain apple")
	}
	// deselect already selected item
	model = &Model{
		Selected: []string{"apple"},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{}) {
		t.Fatal("multiselect should be empty")
	}
	// select multiple items
	model = &Model{
		Selected: []string{"apple"},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 1},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{"apple", "banana"}) {
		t.Fatal("multiselect should contain apple and banana")
	}
	// deselect multiple items
	model = &Model{
		Selected: []string{"apple", "banana"},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model.Viewport.Cursor = 1
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{}) {
		t.Fatal("multiselect should be empty")
	}
	// empty filtered list
	model = &Model{
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{}},
		Viewport: &viewport.Viewport{},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{}) {
		t.Fatal("multiselect should be empty")
	}
	// enter with existing selection
	model = &Model{
		Selected: []string{"apple"},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !slices.Equal(model.Selected, []string{"apple"}) {
		t.Fatal("multiselect should contain apple")
	}
	// enter last selection with existing selection
	model = &Model{
		Selected: []string{"apple"},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 1},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !slices.Equal(model.Selected, []string{"apple"}) {
		t.Fatal("multiselect should contain apple")
	}
	// enter with no selection auto-selects current
	model = &Model{
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !slices.Equal(model.Selected, []string{"apple"}) {
		t.Fatal("multiselect should contain apple")
	}
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("multiselect should quit")
	}
	// enter with empty filtered list
	model = &Model{
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !slices.Equal(model.Selected, []string{}) {
		t.Fatal("multiselect should be empty")
	}
	// ctrl+c
	model = &Model{
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}}},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !slices.Equal(model.Selected, []string{}) {
		t.Fatal("multiselect should be empty")
	}
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("multiselect should interrupt")
	}
	// escape clears search when active
	model = &Model{
		Selected: []string{"apple"},
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
		t.Fatal("multiselect search be inactive")
	}
	if model.Search.Value != "" {
		t.Fatal("multiselect search should be empty")
	}
	if len(model.Search.Filtered) != len(model.Search.Choices) {
		t.Fatal("multiselect choices should be the same as filtered choices")
	}
	// escape quits when search inactive
	model = &Model{
		Selected: []string{"apple"},
		Search:   &search.Search{Active: false},
		Viewport: &viewport.Viewport{},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if len(model.Selected) != 0 {
		t.Fatal("multiselect should be empty")
	}
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("multiselect should quit")
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
		t.Fatal("multiselect search should be active")
	}
	if model.Search.Value != "a" {
		t.Fatal("multiselect search value should be a")
	}
	// select, move down, select again, move down, select again, wrap back and deselect
	model = &Model{
		Selected: []string{},
		Search:   &search.Search{Filtered: []search.Choice{{Id: "apple"}, {Id: "banana"}, {Id: "cherry"}}},
		Viewport: &viewport.Viewport{Cursor: 0, Visible: 3},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{"apple", "banana", "cherry"}) {
		t.Fatal("multiselect should be apple, banana, cherry")
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Equal(model.Selected, []string{"banana", "cherry"}) {
		t.Fatal("multiselect should be banana, cherry")
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
