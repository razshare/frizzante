package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/text"
	flag "github.com/spf13/pflag"
)

func Start(c *cli.Cli) error {
	text.Clrscr()

	if !state.Parsed {
		flag.Parse()
		state.Parsed = true
	}

	if *c.Flags.Help {
		return Help(c)
	}

	if *c.Flags.Version {
		return Version(c)
	}

	if *c.Flags.CreateProject != "" {
		return CreateProject(c, *c.Flags.CreateProject)
	}

	if *c.Flags.Generate != "" {
		return Generate(c, ".", *c.Flags.Generate)
	}

	if *c.Flags.Test {
		return Test(c)
	}

	if *c.Flags.Package {
		return Package(c)
	}

	if *c.Flags.PackageWatch {
		return PackageWatch(c)
	}

	if *c.Flags.Check {
		return Check(c)
	}

	if *c.Flags.Update {
		return Update(c)
	}

	if *c.Flags.Install {
		return Install(c, ".")
	}

	if *c.Flags.Format {
		return Format(c)
	}

	if *c.Flags.Touch {
		return Touch(c)
	}

	if *c.Flags.Clean {
		return Clean(c)
	}

	if *c.Flags.Dev {
		return Dev(c)
	}

	if *c.Flags.Build {
		return Build(c)
	}

	if *c.Flags.Configure {
		return Configure(c, ".")
	}

	if *c.Flags.Welcome {
		return Welcome()
	}

	return Menu(c, false)
}
