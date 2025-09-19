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
	Name string
	Go   string
	Efs  embed.FS
}

type GenerateOptions struct {
	App      string
	Selected string
	Go       string
	Air      string
	Bun      string
	Sqlc     string
	Efs      embed.FS
	Platform platform.Platform
	Auto     bool
}

type TestOptions struct {
	App string
	Go  string
	Bun string
}

type PackageOptions struct {
	App  string
	Bun  string
	Prod bool
}

type PackageWatchOptions struct {
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
	Platform platform.Platform
	Auto     bool
}

type WelcomeOptions struct{}

type NpmOptions struct {
	App   string
	Query string
	Bun   string
}
