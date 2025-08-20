package npmselect

import (
	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
	"time"
)

type Model struct {
	Selected  []string
	Prompt    string
	LastQuery string
	Search    *search.Search
	Viewport  *viewport.Viewport
	Packages  []npm.PackageInfo
	Debounce  time.Duration
	Debouncer *time.Timer
	Error     error
	Loading   bool
	Quitting  bool
	Confirmed bool
}

type SearchResultMsg struct {
	Packages []npm.PackageInfo
	Error    error
}

type DebouncedSearchMsg struct {
	Query string
}
