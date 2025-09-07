package navigate

import (
	"testing"

	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name           string
		filtered       []search.Choice
		initialCursor  int
		initialStart   int
		visibleItems   int
		direction      int
		expectedCursor int
		expectedStart  int
	}{
		{
			name: "move down within viewport",
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered:       []search.Choice{},
			initialCursor:  0,
			initialStart:   0,
			visibleItems:   3,
			direction:      1,
			expectedCursor: 0,
			expectedStart:  0,
		},
		{
			name: "single item list stays at 0",
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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
			filtered: []search.Choice{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search := &search.Search{
				Filtered: tt.filtered,
			}
			viewport := &viewport.Viewport{
				Cursor:  tt.initialCursor,
				Start:   tt.initialStart,
				Visible: tt.visibleItems,
			}

			Apply(search, viewport, tt.direction)

			if viewport.Cursor != tt.expectedCursor {
				t.Errorf("Apply() cursor = %d, want %d", viewport.Cursor, tt.expectedCursor)
			}

			if viewport.Start != tt.expectedStart {
				t.Errorf("Apply() start = %d, want %d", viewport.Start, tt.expectedStart)
			}
		})
	}
}

func TestApplySequence(t *testing.T) {
	filtered := []search.Choice{
		{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}, {Id: "6"},
	}

	search := &search.Search{
		Filtered: filtered,
	}
	viewport := &viewport.Viewport{
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
