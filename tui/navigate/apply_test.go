package navigate

import (
	"testing"

	_search "github.com/razshare/frizzante/tui/search"
	_viewport "github.com/razshare/frizzante/tui/viewport"
)

func TestApply(t *testing.T) {
	var viewport *_viewport.Viewport
	var search *_search.Search

	// move down
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 1)
	if viewport.Cursor != 1 {
		t.Fatal("navigate cursor should be 1")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}

	// move down and scroll down
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 2, Offset: 0, Visible: 3}
	Apply(search, viewport, 1)
	if viewport.Cursor != 3 {
		t.Fatal("navigate cursor should be 3")
	}
	if viewport.Offset != 1 {
		t.Fatal("navigate offset should be 2")
	}

	// move down by 3 and scroll down
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 3)
	if viewport.Cursor != 3 {
		t.Fatal("navigate cursor should be 3")
	}
	if viewport.Offset != 1 {
		t.Fatal("navigate offset should be 2")
	}

	// move up
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 2, Offset: 0, Visible: 3}
	Apply(search, viewport, -1)
	if viewport.Cursor != 1 {
		t.Fatal("navigate cursor should be 1")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}

	// move up and scroll up
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 3, Offset: 3, Visible: 2}
	Apply(search, viewport, -1)
	if viewport.Cursor != 2 {
		t.Fatal("navigate cursor should be 2")
	}
	if viewport.Offset != 2 {
		t.Fatal("navigate offset should be 2")
	}

	// move down and wrap around
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 4, Offset: 0, Visible: 3}
	Apply(search, viewport, 1)
	if viewport.Cursor != 0 {
		t.Fatal("navigate cursor should be 0")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}

	// move up and wrap around
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, -1)
	if viewport.Cursor != 4 {
		t.Fatal("navigate cursor should be 4")
	}
	if viewport.Offset != 2 {
		t.Fatal("navigate offset should be 2")
	}

	// empty list does nothing
	search = &_search.Search{Filtered: []_search.Choice{}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 1)
	if viewport.Cursor != 0 {
		t.Fatal("navigate cursor should be 0")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}

	// single item list stays at 0
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 1)
	if viewport.Cursor != 0 {
		t.Fatal("navigate cursor should be 0")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}

	// move down by 3 and scroll down
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 3)
	if viewport.Cursor != 3 {
		t.Fatal("navigate cursor should be 3")
	}
	if viewport.Offset != 1 {
		t.Fatal("navigate offset should be 1")
	}

	// move up by 3 and scroll up
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, -2)
	if viewport.Cursor != 3 {
		t.Fatal("navigate cursor should be 2")
	}
	if viewport.Offset != 1 {
		t.Fatal("navigate offset should be 1")
	}

	// move down by 7 and scroll down
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, 13)
	if viewport.Cursor != 3 {
		t.Fatal("navigate cursor should be 3")
	}
	if viewport.Offset != 1 {
		t.Fatal("navigate offset should be 1")
	}

	// move up by 7 and scroll up
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 3}
	Apply(search, viewport, -11)
	if viewport.Cursor != 4 {
		t.Fatal("navigate cursor should be 4")
	}
	if viewport.Offset != 2 {
		t.Fatal("navigate offset should be 2")
	}

	// viewport larger than list
	search = &_search.Search{Filtered: []_search.Choice{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}}}
	viewport = &_viewport.Viewport{Cursor: 0, Offset: 0, Visible: 12}
	Apply(search, viewport, 1)
	if viewport.Cursor != 1 {
		t.Fatal("navigate cursor should be 1")
	}
	if viewport.Offset != 0 {
		t.Fatal("navigate offset should be 0")
	}
}
