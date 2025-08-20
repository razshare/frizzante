package npmselect

import (
	"time"

	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Objects []struct {
		Package PackageInfo `json:"package"`
	} `json:"objects"`
}

type Model struct {
	Search        *search.Search
	Viewport      *viewport.Viewport
	Packages      []PackageInfo
	Selected      []string
	Loading       bool
	Error         error
	LastQuery     string
	DebounceTimer *time.Timer
	Quitting      bool
	Confirmed     bool
}

type SearchResultMsg struct {
	Packages []PackageInfo
	Error    error
}

type DebouncedSearchMsg struct {
	Query string
}