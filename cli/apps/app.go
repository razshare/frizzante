package apps

import "embed"

type App struct {
	Efs            embed.FS
	Go             *string
	Air            *string
	Bun            *string
	Sqlc           *string
	SqlcYaml       *string
	Tags           *string
	Database       *string
	DatabaseType   *string
	Create         *string
	Add            *string
	Context        *string
	Migrate        *string
	Generate       *string
	Strict         *bool
	Dev            *bool
	Configure      *bool
	Install        *bool
	Update         *bool
	Build          *bool
	Asm            *bool
	Package        *bool
	PackageWatch   *bool
	Check          *bool
	Format         *bool
	Touch          *bool
	Clean          *bool
	Reset          *bool
	Clear          *bool
	LockJsPackages *bool
	Test           *bool
	Welcome        *bool
	Help           *bool
	Version        *bool
}
