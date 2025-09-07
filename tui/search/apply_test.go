package search

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name          string
		choices       []Choice
		searchInput   string
		expectedCount int
		expectedFirst string
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search := &Search{
				Choices:  tt.choices,
				Filtered: tt.choices,
				Input:    textinput.New(),
			}
			viewport := &viewport.Viewport{}

			search.Input.SetValue(tt.searchInput)
			Filter(search, viewport)

			if len(search.Filtered) != tt.expectedCount {
				t.Errorf("Filter() resulted in %d items, want %d", len(search.Filtered), tt.expectedCount)
			}

			if tt.expectedCount > 0 && search.Filtered[0].Id != tt.expectedFirst {
				t.Errorf("Filter() first item = %s, want %s", search.Filtered[0].Id, tt.expectedFirst)
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

	viewport := &viewport.Viewport{
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
		viewport := &viewport.Viewport{}

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
