package search

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	_viewport "github.com/razshare/frizzante/tui/viewport"
)

func TestFilter(t *testing.T) {
	type TestData struct {
		name          string
		choices       []Choice
		searchInput   string
		expectedCount int
		expectedFirst string
	}

	data := []TestData{
		{
			name: "empty search returns all choices",
			choices: []Choice{
				{Id: "apple", Description: "A fruit"},
				{Id: "banana", Description: "Another fruit"},
				{Id: "cherry", Description: "Small fruit"},
			},
			searchInput:   "",
			expectedCount: 3,
			expectedFirst: "apple",
		},
		{
			name: "case insensitive search",
			choices: []Choice{
				{Id: "Apple", Description: "A fruit"},
				{Id: "Banana", Description: "Another fruit"},
				{Id: "Cherry", Description: "Small fruit"},
			},
			searchInput:   "APPLE",
			expectedCount: 1,
			expectedFirst: "Apple",
		},
		{
			name: "partial match",
			choices: []Choice{
				{Id: "apple", Description: "A fruit"},
				{Id: "pineapple", Description: "Tropical fruit"},
				{Id: "banana", Description: "Another fruit"},
			},
			searchInput:   "app",
			expectedCount: 2,
			expectedFirst: "apple",
		},
		{
			name: "no matches",
			choices: []Choice{
				{Id: "apple", Description: "A fruit"},
				{Id: "banana", Description: "Another fruit"},
			},
			searchInput:   "xyz",
			expectedCount: 0,
		},
		{
			name: "substring match",
			choices: []Choice{
				{Id: "development", Description: "Dev mode"},
				{Id: "production", Description: "Prod mode"},
				{Id: "developer", Description: "Dev user"},
			},
			searchInput:   "dev",
			expectedCount: 2,
			expectedFirst: "development",
		},
		{
			name:          "empty choices",
			choices:       []Choice{},
			searchInput:   "test",
			expectedCount: 0,
		},
		{
			name: "whitespace in search",
			choices: []Choice{
				{Id: "hello world", Description: "Greeting"},
				{Id: "world peace", Description: "Goal"},
				{Id: "hello", Description: "Simple greeting"},
			},
			searchInput:   "hello",
			expectedCount: 2,
			expectedFirst: "hello world",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			search := &Search{
				Choices:  d.choices,
				Filtered: d.choices,
				Input:    textinput.New(),
			}

			viewport := &_viewport.Viewport{}

			search.Input.SetValue(d.searchInput)
			Filter(search, viewport)

			if len(search.Filtered) != d.expectedCount {
				t.Errorf("Filter() resulted in %d items, want %d", len(search.Filtered), d.expectedCount)
			}

			if d.expectedCount > 0 && search.Filtered[0].Id != d.expectedFirst {
				t.Errorf("Filter() first item = %s, want %s", search.Filtered[0].Id, d.expectedFirst)
			}

			if viewport.Cursor != 0 {
				t.Errorf("Filter() cursor = %d, want 0", viewport.Cursor)
			}

			if viewport.Start != 0 {
				t.Errorf("Filter() start = %d, want 0", viewport.Start)
			}
		})
	}
}

func TestReset(t *testing.T) {
	choices := []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "banana", Description: "Another fruit"},
		{Id: "cherry", Description: "Small fruit"},
	}

	search := &Search{
		Active:   true,
		Choices:  choices,
		Filtered: []Choice{{Id: "apple", Description: "A fruit"}},
		Input:    textinput.New(),
	}

	viewport := &_viewport.Viewport{
		Cursor: 5,
		Start:  2,
	}

	search.Input.SetValue("test")

	Reset(search, viewport)

	if search.Active {
		t.Error("Reset() Active should be false")
	}

	if search.Input.Value() != "" {
		t.Errorf("Reset() Input.Value() = %s, want empty", search.Input.Value())
	}

	if len(search.Filtered) != len(search.Choices) {
		t.Errorf("Reset() Filtered length = %d, want %d", len(search.Filtered), len(search.Choices))
	}

	for i, choice := range search.Filtered {
		if choice.Id != search.Choices[i].Id {
			t.Errorf("Reset() Filtered[%d] = %s, want %s", i, choice.Id, search.Choices[i].Id)
		}
	}

	if viewport.Cursor != 0 {
		t.Errorf("Reset() Cursor = %d, want 0", viewport.Cursor)
	}

	if viewport.Start != 0 {
		t.Errorf("Reset() Start = %d, want 0", viewport.Start)
	}
}

func TestSearchStateTransitions(t *testing.T) {
	choices := []Choice{
		{Id: "apple", Description: "A fruit"},
		{Id: "banana", Description: "Another fruit"},
		{Id: "application", Description: "Software"},
	}

	t.Run("activate search", func(t *testing.T) {
		search := &Search{
			Active:   false,
			Choices:  choices,
			Filtered: choices,
			Input:    textinput.New(),
		}

		search.Active = true
		search.Input.Focus()

		if !search.Active {
			t.Error("Search should be active")
		}
	})

	t.Run("search then reset flow", func(t *testing.T) {
		search := &Search{
			Active:   false,
			Choices:  choices,
			Filtered: choices,
			Input:    textinput.New(),
		}
		viewport := &_viewport.Viewport{}

		search.Active = true
		search.Input.SetValue("app")
		Filter(search, viewport)

		if len(search.Filtered) != 2 {
			t.Errorf("After filter, expected 2 results, got %d", len(search.Filtered))
		}

		Reset(search, viewport)

		if search.Active {
			t.Error("After reset, search should be inactive")
		}

		if len(search.Filtered) != len(choices) {
			t.Errorf("After reset, expected all choices, got %d", len(search.Filtered))
		}
	})
}
