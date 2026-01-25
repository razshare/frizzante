package spinners

import (
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate(t *testing.T) {
	var model *Model
	var cmd tea.Cmd
	// hard interrupt with ctrl+c
	model = &Model{Message: "Loading...", Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should interrupt")
	}
	// escape does nothing
	model = &Model{Message: "Loading...", Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd != nil {
		t.Fatal("spinner should ignore esc key")
	}
	// tick
	model = &Model{Message: "Loading...", Spinner: spinner.New()}
	cmd = model.Init()
	if _, ok := reflect.TypeAssert[tea.Msg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should tick")
	}
	// view
	model = &Model{Message: "Test Message", Spinner: spinner.New()}
	view := model.View()
	if view == "" {
		t.Fatal("spinner view should not be empty")
	}
	if !strings.Contains(view, "Test Message") {
		t.Fatal("view view should contain Test Message")
	}
}
