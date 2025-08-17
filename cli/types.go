package cli

import (
	"embed"
	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/platform"
)

type FeatureCopyInstruction struct {
	FeatureName          string
	OriginDirectory      string
	DestinationDirectory string
}

type Menu map[string]func() (func(), error)

type Cli struct {
	App          *string
	Help         *bool
	Version      *bool
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
	Clean        *bool
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
	Menu         Menu
}

type NextOptions struct {
	Go       string
	Air      string
	Bun      string
	Sqlc     string
	Platform platform.Platform
	Action   func() action.Type
	Counter  uint64
}

type SelectOptions struct {
	Go          string
	Air         string
	Bun         string
	Sqlc        string
	Platform    platform.Platform
	Action      func() action.Type
	NextCounter uint64
}
