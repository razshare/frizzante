package apps

import "embed"

type App struct {
	Add              *string
	Help             *bool
	Version          *bool
	Reset            *bool
	CreateProject    *string
	Generate         *string
	Migrate          *string
	Test             *bool
	Package          *bool
	PackageWatch     *bool
	Check            *bool
	Update           *bool
	Install          *bool
	Format           *bool
	Touch            *bool
	CleanProject     *bool
	Dev              *bool
	Build            *bool
	Configure        *bool
	Yes              *bool
	Welcome          *bool
	Clear            *bool
	AssemblyExplorer *bool
	LockPackages     *bool
	Go               *string
	Air              *string
	Bun              *string
	Sqlc             *string
	SqlcYaml         *string
	Tags             *string
	Database         *string
	Efs              embed.FS
}
