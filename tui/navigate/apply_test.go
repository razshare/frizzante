package navigate

import (
	"testing"

	_search "github.com/razshare/frizzante/tui/search"
	_viewport "github.com/razshare/frizzante/tui/viewport"
)

func TestApply(t *testing.T) {
	type TestData struct {
		name           string
		filtered       []_search.Choice
		initialCursor  int
		initialStart   int
		visibleItems   int
		direction      int
		expectedCursor int
		expectedStart  int
	}

	data := []TestData{
		{
			name: "move down within viewport",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 1,
			expectedStart:  0,
		},
		{
			name: "move down triggers scroll",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  2,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 3,
			expectedStart:  1,
		},
		{
			name: "move up within viewport",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  2,
			initialStart:   0,
			visibleItems:   3,
			direction:      -1,
			expectedCursor: 1,
			expectedStart:  0,
		},
		{
			name: "move up triggers scroll",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  3,
			initialStart:   3,
			visibleItems:   2,
			direction:      -1,
			expectedCursor: 2,
			expectedStart:  2,
		},
		{
			name: "wrap around from bottom to top",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"},
			},
			initialCursor:  2,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 0,
			expectedStart:  0,
		},
		{
			name: "wrap around from top to bottom",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   2,
			direction:      -1,
			expectedCursor: 2,
			expectedStart:  1,
		},
		{
			name:           "empty list does nothing",
			filtered:       []_search.Choice{},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 0,
			expectedStart:  0,
		},
		{
			name: "single item list stays at 0",
			filtered: []_search.Choice{
				{Id: "1"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 0,
			expectedStart:  0,
		},
		{
			name: "move multiple steps down",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   3,
			direction:      3,
			expectedCursor: 3,
			expectedStart:  1,
		},
		{
			name: "move multiple steps up",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"},
			},
			initialCursor:  4,
			initialStart:   2,
			visibleItems:   3,
			direction:      -3,
			expectedCursor: 1,
			expectedStart:  1,
		},
		{
			name: "large positive direction wraps correctly",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   2,
			direction:      7,
			expectedCursor: 1,
			expectedStart:  0,
		},
		{
			name: "large negative direction wraps correctly",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"}, {Id: "3"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   2,
			direction:      -7,
			expectedCursor: 2,
			expectedStart:  1,
		},
		{
			name: "viewport larger than list",
			filtered: []_search.Choice{
				{Id: "1"}, {Id: "2"},
			},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   5,
			direction:      1,
			expectedCursor: 1,
			expectedStart:  0,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			search := &_search.Search{
				Filtered: d.filtered,
			}
			viewport := &_viewport.Viewport{
				Cursor:  d.initialCursor,
				Start:   d.initialStart,
				Visible: d.visibleItems,
			}

			Apply(search, viewport, d.direction)

			if viewport.Cursor != d.expectedCursor {
				t.Errorf("Apply() cursor = %d, want %d", viewport.Cursor, d.expectedCursor)
			}

			if viewport.Start != d.expectedStart {
				t.Errorf("Apply() start = %d, want %d", viewport.Start, d.expectedStart)
			}
		})
	}
}

func TestApplySequence(t *testing.T) {
	filtered := []_search.Choice{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}, {Id: "6"},
	}

	search := &_search.Search{
		Filtered: filtered,
	}
	viewport := &_viewport.Viewport{
		Cursor:  0,
		Start:   0,
		Visible: 3,
	}

	movements := []struct {
		direction      int
		expectedCursor int
		expectedStart  int
	}{
		{1, 1, 0},
		{1, 2, 0},
		{1, 3, 1},
		{1, 4, 2},
		{1, 5, 3},
		{1, 0, 0},
		{-1, 5, 3},
		{-1, 4, 3},
		{-1, 3, 3},
		{-1, 2, 2},
		{-1, 1, 1},
		{-1, 0, 0},
	}

	for i, move := range movements {
		Apply(search, viewport, move.direction)
		if viewport.Cursor != move.expectedCursor {
			t.Errorf("Step %d: cursor = %d, want %d", i, viewport.Cursor, move.expectedCursor)
		}
		if viewport.Start != move.expectedStart {
			t.Errorf("Step %d: start = %d, want %d", i, viewport.Start, move.expectedStart)
		}
	}
}
