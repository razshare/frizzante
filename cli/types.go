package cli

import "embed"

type FeatureAddEvents struct {
	ConfirmOverwrite            func(name string) bool
	ConfirmAddMissingDependency func(feature string, dependency string) bool
}

type InstallEvents struct {
	ConfirmOverwrite func(name string) bool
}

type FeatureCopyInstruction struct {
	From string
	To   string
}

type Cli struct {
	Efs    embed.FS
	Parsed bool
}
