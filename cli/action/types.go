package action

import (
	"embed"
	"time"

	"github.com/razshare/frizzante/platform"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type HelpOptions struct{}

type ClearOptions struct{}

type ResetOptions struct{}

type VersionOptions struct {
	Efs embed.FS
}

type CreateProjectOptions struct {
	Project string
}

type GenerateOptions struct {
	App      string
	Selected string
	Auto     bool
	Go       string
	Air      string
	Bun      string
	Sqlc     string
	Efs      embed.FS
	Platform platform.Platform
}

type TestOptions struct {
	App string
	Go  string
	Bun string
}

type PkgOptions struct {
	App string
	Bun string
}

type PkgWatchOptions struct {
	App string
	Bun string
}

type CheckOptions struct {
	App string
	Bun string
}

type InstallOptions struct {
	App string
	Go  string
	Bun string
}

type UpdateOptions struct {
	App string
	Go  string
	Bun string
}

type FormatOptions struct {
	App string
	Go  string
	Bun string
}

type TouchOptions struct {
	App string
}

type CleanProjectOptions struct {
	App string
	Go  string
}

type DevOptions struct {
	App string
	Go  string
	Air string
	Bun string
}

type BuildOptions struct {
	App      string
	Go       string
	Bun      string
	Platform platform.Platform
}

type ConfigOptions struct {
	App      string
	Clear    bool
	Auto     bool
	Go       string
	Air      string
	Bun      string
	Sqlc     string
	Generate string
	Efs      embed.FS
	Platform platform.Platform
}

type WelcomeOptions struct{}

type NpmOptions struct {
	Query string
}

type NpmPackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type NpmSearchResponse struct {
	Objects []struct {
		Package NpmPackageInfo `json:"package"`
	} `json:"objects"`
}

type NpmSearchModel struct {
	Search        *search.Search
	Viewport      *viewport.Viewport
	Packages      []NpmPackageInfo
	Selected      []string
	Loading       bool
	Error         error
	LastQuery     string
	DebounceTimer *time.Timer
	Quitting      bool
	Confirmed     bool
}

type SearchResultMsg struct {
	Packages []NpmPackageInfo
	Error    error
}

type DebouncedSearchMsg struct {
	Query string
}