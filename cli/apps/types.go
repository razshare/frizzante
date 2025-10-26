package apps

import "embed"

type App struct {
	Add              *bool
	Help             *bool
	Version          *bool
	Reset            *bool
	CreateProject    *bool
	Generate         *bool
	Migrate          *bool
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
	Platform         *string
	Go               *string
	Air              *string
	Bun              *string
	Sqlc             *string
	SqlcYaml         *string
	Tags             *string
	Database         *string
	Value            *string
	Efs              embed.FS
}
