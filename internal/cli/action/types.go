package action

import (
	"embed"

	"github.com/razshare/frizzante/internal/platform"
)

type HelpOptions struct{}

type ClearOptions struct{}

type ResetOptions struct{}

type VersionOptions struct {
	Efs embed.FS
}

type CreateProjectOptions struct {
	Name string
	Efs  embed.FS
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

type PacakgeOptions struct {
	App string
	Bun string
}

type PacakgeWatchOptions struct {
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

type ConfigureOptions struct {
	App      string
	Go       string
	Air      string
	Bun      string
	Auto     bool
	Platform platform.Platform
}

type WelcomeOptions struct{}

type NpmOptions struct {
	App   string
	Query string
	Bun   string
}
