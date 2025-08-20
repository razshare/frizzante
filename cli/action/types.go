package action

import (
	"embed"

	"github.com/razshare/frizzante/platform"
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
	App   string
	Query string
}

