package select_npm_packages

import (
	"time"

	"github.com/razshare/frizzante/v2/cli/npm"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
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
