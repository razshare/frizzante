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
	Value string
	Go    string
	Air   string
	Bun   string
	Efs   embed.FS
}

type GenerateOptions struct {
	Value    string
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
}

type TestOptions struct {
	Go  string
	Bun string
}

type PackageOptions struct {
	Bun  string
	Prod bool
}

type PackageWatchOptions struct {
	Bun      string
	Callback func()
}

type CheckOptions struct {
	Bun string
}

type InstallOptions struct {
	Go  string
	Bun string
}

type UpdateOptions struct {
	Go  string
	Bun string
}

type FormatOptions struct {
	Go  string
	Bun string
}

type TouchOptions struct {
}

type CleanProjectOptions struct {
	Go string
}

type DevOptions struct {
	Go   string
	Air  string
	Bun  string
	Efs  embed.FS
	Tags []string
}

type BuildOptions struct {
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
}

type AssemblyExplorerOptions struct {
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
	Auto     bool
}

type ConfigureOptions struct {
	Go       string
	Air      string
	Bun      string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
}

type WelcomeOptions struct{}

type NpmOptions struct {
	Value string
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
