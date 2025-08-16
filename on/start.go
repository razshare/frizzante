package on

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/text"
	flag "github.com/spf13/pflag"
)

func Start(efs embed.FS) error {
	text.Clrscr()

	if !state.Parsed {
		flag.Parse()
		state.Parsed = true
	}

	if *state.Help {
		return Help()
	}

	if *state.Version {
		return Version(efs)
	}

	if *state.CreateProject != "" {
		return CreateProject(efs, *state.CreateProject)
	}

	if *state.Generate != "" {
		return Generate(efs, ".", *state.Generate)
	}

	if *state.Test {
		return Test()
	}

	if *state.Package {
		return Package()
	}

	if *state.PackageWatch {
		return PackageWatch()
	}

	if *state.Check {
		return Check()
	}

	if *state.Update {
		return Update()
	}

	if *state.Install {
		return Install(".")
	}

	if *state.Format {
		return Format()
	}

	if *state.Touch {
		return Touch()
	}

	if *state.Clean {
		return Clean()
	}

	if *state.Dev {
		return Dev()
	}

	if *state.Build {
		return Build()
	}

	if *state.Configure {
		return Configure(efs, ".")
	}

	if *state.Welcome {
		return Welcome()
	}

	return Menu(efs, false)
}
