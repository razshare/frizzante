package actions

import (
	"database/sql"
	"embed"

	"github.com/razshare/frizzante/platforms"
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
	SqlcYaml string
	Database string
	Tags     []string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
	Active   bool
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
	App      string
	Bun      string
	Callback func()
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
	App  string
	Go   string
	Air  string
	Bun  string
	Efs  embed.FS
	Tags []string
}

type BuildOptions struct {
	App      string
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
}

type AssemblyExplorerOptions struct {
	App      string
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
	Auto     bool
}

type ConfigureOptions struct {
	App      string
	Go       string
	Air      string
	Bun      string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
}

type WelcomeOptions struct{}

type NpmOptions struct {
	App   string
	Query string
	Bun   string
}

type MigrateOptions struct {
	Offset   string
	Target   string
	Sqlc     string
	SqlcYaml string
	Database *sql.DB
	Platform platforms.Platform
	Auto     bool
}
