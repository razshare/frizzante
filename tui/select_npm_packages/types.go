package select_npm_packages

import (
	"time"

	"github.com/razshare/frizzante/cli/npm"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type Model struct {
	Packages  []npm.PackageInfo
	Selected  []string
	Prompt    string
	LastQuery string
	Error     error
	Search    *search.Search
	Viewport  *viewport.Viewport
	Debounce  time.Duration
	Debouncer *time.Timer
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
