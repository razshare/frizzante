package cli

import "embed"

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}

type Flags struct {
	App           *string
	Help          *bool
	Version       *bool
	CreateProject *string
	Generate      *string
	Test          *bool
	Package       *bool
	PackageWatch  *bool
	Check         *bool
	Update        *bool
	Install       *bool
	Format        *bool
	Touch         *bool
	Clean         *bool
	Dev           *bool
	Build         *bool
	Configure     *bool
	Platform      *string
	Yes           *bool
	Go            *string
	Air           *string
	Bun           *string
	Sqlc          *string
	Welcome       *bool
}

type Cli struct {
	Efs   embed.FS
	Flags *Flags
}
