package npmselect

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestSearch(t *testing.T) {
	var model *Model
	var cmd tea.Cmd

	// debouncing
	model = &Model{
		Selected:  []string{},
		Search:    &search.Search{Active: false, Input: textinput.New()},
		Viewport:  &viewport.Viewport{},
		Debouncer: time.NewTimer(100 * time.Millisecond),
		Debounce:  100 * time.Millisecond,
	}
	model.Debouncer.Stop()
	model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if !model.Search.Active {
		t.Fatal("search should be activated")
	}
	if model.Search.Input.Value() != "r" {
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
			Input: textinput.New(),
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
		Search:   &search.Search{Input: textinput.New()},
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
		Search:   &search.Search{Filtered: []search.Choice{}, Input: textinput.New()},
		Viewport: &viewport.Viewport{},
	}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok || !model.Confirmed {
		t.Fatal("search should confirm")
	}
}
