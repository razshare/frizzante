package app

import "embed"

type App struct {
	App          *string
	Help         *bool
	Version      *bool
	Reset        *bool
	Project      *string
	Generate     *string
	Test         *bool
	Package      *bool
	PackageWatch *bool
	Check        *bool
	Update       *bool
	Install      *bool
	Format       *bool
	Touch        *bool
	CleanProject *bool
	Dev          *bool
	Build        *bool
	Configure    *bool
	Platform     *string
	Yes          *bool
	Go           *string
	Air          *string
	Bun          *string
	Sqlc         *string
	Welcome      *bool
	Clear        *bool
	Efs          embed.FS
}
