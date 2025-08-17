package action

import (
	"embed"
	"github.com/razshare/frizzante/platform"
)

type Type int

const (
	TypeMenu     Type = 0
	TypeVersion  Type = 1
	TypeProject  Type = 2
	TypeGenerate Type = 3
	TypeTest     Type = 4
	TypePkg      Type = 5
	TypePkgWatch Type = 6
	TypeCheck    Type = 7
	TypeInstall  Type = 8
	TypeUpdate   Type = 9
	TypeFormat   Type = 10
	TypeTouch    Type = 11
	TypeClean    Type = 12
	TypeDev      Type = 13
	TypeBuild    Type = 14
	TypeConfig   Type = 15
	TypeWelcome  Type = 16
	TypeHelp     Type = 17
)

type HelpOptions struct{}

type VersionOptions struct {
	Efs embed.FS
}

type CreateOptions struct {
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

type CleanOptions struct {
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
	Platform platform.Platform
}

type WelcomeOptions struct{}
