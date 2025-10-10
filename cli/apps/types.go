package apps

import "embed"

type App struct {
	Add              *string
	App              *string
	Help             *bool
	Version          *bool
	Reset            *bool
	CreateProject    *string
	Generate         *string
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
	Platform         *string
	Yes              *bool
	Go               *string
	Air              *string
	Bun              *string
	Sqlc             *string
	Welcome          *bool
	Clear            *bool
	Tags             *string
	AssemblyExplorer *bool
	Efs              embed.FS
}
