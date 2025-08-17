package cli

import (
	"fmt"
	"github.com/razshare/frizzante/cli/action"
	"github.com/razshare/frizzante/tui/config"
)

func Next(c *Cli, o NextOptions) error {
	switch o.Action() {
	case action.TypeHelp:
		return action.Help(action.HelpOptions{})
	case action.TypeVersion:
		return action.Version(action.VersionOptions{
			Efs: c.Efs,
		})
	case action.TypeProject:
		return action.Create(action.CreateOptions{
			Project: *c.Project,
		})
	case action.TypeGenerate:
		return action.Generate(action.GenerateOptions{
			App:      *c.App,
			Selected: *c.Generate,
			Auto:     *c.Yes,
			Platform: o.Platform,
			Go:       o.Go,
			Air:      o.Air,
			Bun:      o.Bun,
			Sqlc:     o.Sqlc,
		})
	case action.TypeTest:
		return action.Test(action.TestOptions{
			App: *c.App,
			Go:  o.Go,
			Bun: o.Bun,
		})
	case action.TypePkg:
		return action.Pkg(action.PkgOptions{
			App: *c.App,
			Bun: o.Bun,
		})
	case action.TypePkgWatch:
		return action.PkgWatch(action.PkgWatchOptions{
			App: *c.App,
			Bun: o.Bun,
		})
	case action.TypeCheck:
		return action.Check(action.CheckOptions{
			App: *c.App,
			Bun: o.Bun,
		})
	case action.TypeInstall:
		return action.Install(action.InstallOptions{
			App: *c.App,
			Go:  o.Go,
			Bun: o.Bun,
		})
	case action.TypeUpdate:
		return action.Update(action.UpdateOptions{
			App: *c.App,
			Go:  o.Go,
			Bun: o.Bun,
		})
	case action.TypeFormat:
		return action.Format(action.FormatOptions{
			App: *c.App,
			Go:  o.Go,
			Bun: o.Bun,
		})
	case action.TypeTouch:
		return action.Touch(action.TouchOptions{
			App: *c.App,
		})
	case action.TypeClean:
		return action.Clean(action.CleanOptions{
			App: *c.App,
			Go:  o.Go,
		})
	case action.TypeDev:
		return action.Dev(action.DevOptions{
			App: *c.App,
			Go:  o.Go,
			Air: o.Air,
			Bun: o.Bun,
		})
	case action.TypeBuild:
		return action.Build(action.BuildOptions{
			App:      *c.App,
			Platform: o.Platform,
			Go:       o.Go,
			Bun:      o.Bun,
		})
	case action.TypeConfig:
		return action.Config(action.ConfigOptions{
			App:      *c.App,
			Generate: *c.Generate,
			Clear:    *c.Clear,
			Auto:     *c.Yes,
			Platform: o.Platform,
			Go:       o.Go,
			Air:      o.Air,
			Bun:      o.Bun,
			Sqlc:     o.Sqlc,
		})
	case action.TypeWelcome:
		return action.Welcome(action.WelcomeOptions{})
	default:
		if o.Counter == 0 {
			logo, err := c.Efs.ReadFile("clilogo.txt")
			if err != nil {
				return err
			}
			fmt.Println(config.Styles.BigText.Render(string(logo)))
			o.Counter++
		}
		return Select(c, SelectOptions{
			Go:          o.Go,
			Air:         o.Air,
			Bun:         o.Bun,
			Sqlc:        o.Sqlc,
			Platform:    o.Platform,
			Action:      o.Action,
			NextCounter: o.Counter,
		})
	}
}
