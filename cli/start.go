package cli

import (
	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/cli/path"
)

func Start(c *Cli) error {
	plat, err := Platform(c)
	if err != nil {
		return err
	}

	_go, err := path.Go(*c.Go)
	if err != nil {
		return err
	}

	air, err := path.Air(*c.Air)
	if err != nil {
		return err
	}

	bun, err := path.Bun(*c.Bun)
	if err != nil {
		return err
	}

	sqlc, err := path.Sqlc(*c.Sqlc)
	if err != nil {
		return err
	}

	return Next(c, NextOptions{
		Go:       _go,
		Air:      air,
		Bun:      bun,
		Sqlc:     sqlc,
		Platform: plat,
		Action: func() action.Type {
			if *c.Help {
				return action.TypeHelp
			}

			if *c.Version {
				return action.TypeVersion
			}

			if *c.Reset {
				return action.TypeReset
			}

			if *c.Project != "" {
				return action.TypeProject
			}

			if *c.Generate != "" {
				return action.TypeGenerate
			}

			if *c.Test {
				return action.TypeTest
			}

			if *c.Package {
				return action.TypePkg
			}

			if *c.PackageWatch {
				return action.TypePkgWatch
			}

			if *c.Check {
				return action.TypeCheck
			}

			if *c.Install {
				return action.TypeInstall
			}

			if *c.Update {
				return action.TypeUpdate
			}

			if *c.Format {
				return action.TypeFormat
			}

			if *c.Touch {
				return action.TypeTouch
			}

			if *c.Clean {
				return action.TypeClean
			}

			if *c.Dev {
				return action.TypeDev
			}

			if *c.Build {
				return action.TypeBuild
			}

			if *c.Configure {
				return action.TypeConfig
			}

			if *c.Welcome {
				return action.TypeWelcome
			}

			return action.TypeMenu
		},
	})
}
