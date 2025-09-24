package spinner

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

	// soft interrupt with ctrl+c
	model = &Model{SoftInterrupt: true, Message: "Loading...", Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should quit")
	}

	// hard interrupt with ctrl+c
	model = &Model{Message: "Loading...", Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should interrupt")
	}

	// escape quits
	model = &Model{Message: "Loading...", SoftInterrupt: true, Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should quit")
	}

	// escape quits even with hard interrupt
	model = &Model{Message: "Loading...", Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should quit")
	}

	// soft to hard interrupt transition
	model = &Model{Message: "Loading...", SoftInterrupt: true, Spinner: spinner.New()}
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.QuitMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should quit")
	}
	model.SoftInterrupt = false
	_, cmd = model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if _, ok := reflect.TypeAssert[tea.InterruptMsg](reflect.ValueOf(cmd())); !ok {
		t.Fatal("spinner should interrupt")
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
