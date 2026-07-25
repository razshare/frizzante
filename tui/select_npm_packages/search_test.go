package select_npm_packages

import (
	"reflect"
	"slices"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
)

func TestSearch(t *testing.T) {
	var model *Model
	var cmd tea.Cmd
	// debouncing
	model = &Model{
		Selected:  []string{},
		Search:    &search.Search{Active: false},
		Viewport:  &viewport.Viewport{},
		Debouncer: time.NewTimer(100 * time.Millisecond),
		Debounce:  100 * time.Millisecond,
	}
	model.Debouncer.Stop()
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if !model.Search.Active {
		t.Fatal("search should be activated")
	}
	if model.Search.Value != "r" {
		t.Fatal("search input should be r")
	}
	if model.LastQuery != "r" {
		t.Fatal("search input last query should be r")
	}
	// space selects
	model = &Model{
		Selected: []string{},
		Search: &search.Search{
			Filtered: []search.Choice{
				{Id: "react@18.2.0"},
				{Id: "vue@3.3.0"},
				{Id: "angular@16.0.0"},
			},
		},
		Viewport: &viewport.Viewport{Cursor: 0},
	}
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !slices.Contains(model.Selected, "react@18.2.0") {
		t.Fatal("search should select react@18.2.0")
	}
	model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if len(model.Selected) != 2 {
		t.Fatal("search should select 2 items")
	}
	if !slices.Contains(model.Selected, "vue@3.3.0") {
		t.Fatal("search should select react@18.2.0 and vue@3.3.0")
	}
	model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model.Update(tea.KeyMsg{Type: tea.KeySpace})
	if slices.Contains(model.Selected, "react@18.2.0") {
		t.Fatal("search should not select react@18.2.0")
	}
	// esc quits
	model = &Model{
		Selected: []string{"react@18.2.0"},
		Search:   &search.Search{},
		Viewport: &viewport.Viewport{},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok || !model.Quitting {
		t.Fatal("search should quit")
	}
	if len(model.Selected) != 0 {
		t.Fatal("search should not select anything")
	}
	// enter confirms
	model = &Model{
		Selected: []string{"react@18.2.0"},
		Search:   &search.Search{Filtered: []search.Choice{}},
		Viewport: &viewport.Viewport{},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok || !model.Confirmed {
		t.Fatal("search should confirm")
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
