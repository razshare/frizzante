package search

import (
	"testing"

	viewport_ "github.com/razshare/frizzante/tui/viewport"
)

func TestApply(t *testing.T) {
	var search *Search
	var choices []Choice
	var viewport *viewport_.Viewport

	// empty search returns all choices
	choices = []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "banana", Description: "Another fruit"},
		{Id: "cherry", Description: "Small fruit"},
	}
	search = &Search{Choices: choices, Filtered: choices}
	viewport = &viewport_.Viewport{}
	search.Value = ""
	Filter(search, viewport)
	if len(search.Filtered) != 3 {
		t.Fatal("search filter should contain 3 items")
	}
	if viewport.Cursor != 0 {
		t.Fatal("search cursor should be 0")
	}
	if viewport.Offset != 0 {
		t.Fatal("search offset should be 0")
	}

	// case-insensitive search
	choices = []Choice{
		{Id: "Apple", Description: "A fruit"},
		{Id: "Banana", Description: "Another fruit"},
		{Id: "Cherry", Description: "Small fruit"},
	}
	search = &Search{Choices: choices, Filtered: choices}
	viewport = &viewport_.Viewport{}
	search.Value = "APPLE"
	Filter(search, viewport)
	if len(search.Filtered) != 1 {
		t.Fatal("search filter should contain 1 item")
	}
	if search.Filtered[0].Id != "Apple" {
		t.Fatal("search filter should contain Apple")
	}

	// partial match
	choices = []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "pineapple", Description: "Tropical fruit"},
		{Id: "banana", Description: "Another fruit"},
	}
	search = &Search{Choices: choices, Filtered: choices}
	viewport = &viewport_.Viewport{}
	search.Value = "app"
	Filter(search, viewport)
	if len(search.Filtered) != 2 {
		t.Fatal("search filter should contain 2 items")
	}
	if search.Filtered[0].Id != "apple" || search.Filtered[1].Id != "pineapple" {
		t.Fatal("search filter should contain apple (1st) and pineapple (2nd)")
	}

	// no matches
	choices = []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "banana", Description: "Another fruit"},
	}
	search = &Search{Choices: choices, Filtered: choices}
	viewport = &viewport_.Viewport{}
	search.Value = "xyz"
	Filter(search, viewport)
	if len(search.Filtered) != 0 {
		t.Fatal("search filter should be empty")
	}

	// empty choices
	choices = []Choice{}
	search = &Search{Choices: choices, Filtered: choices}
	viewport = &viewport_.Viewport{}
	search.Value = "test"
	Filter(search, viewport)
	if len(search.Filtered) != 0 {
		t.Fatal("search filter should be empty")
	}

	// reset then reset
	choices = []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "banana", Description: "Another fruit"},
		{Id: "cherry", Description: "Small fruit"},
	}
	search = &Search{Choices: choices, Filtered: choices}
	search.Value = "apple"
	Filter(search, viewport)
	if len(search.Filtered) != 1 {
		t.Fatal("search filter should contain 1 item")
	}
	if search.Filtered[0].Id != "apple" {
		t.Fatal("search filter should contain apple")
	}
	Reset(search, viewport)
	if search.Active {
		t.Fatal("search should not be active")
	}
	if len(search.Filtered) != len(search.Choices) {
		t.Fatal("search filter should empty")
	}
	if search.Value != "" {
		t.Fatal("search input should empty")
	}
	if viewport.Cursor != 0 {
		t.Fatal("search cursor should be 0")
	}
	if viewport.Offset != 0 {
		t.Fatal("search offset should be 0")
	}
}
