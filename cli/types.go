package cli

import "embed"

type FeatureAddEvents struct {
	ConfirmFeatureOverwrite     func(feature string) bool
	ConfirmAddMissingDependency func(feature string, dependency string) bool
}

type Cli struct {
	Efs embed.FS
}
