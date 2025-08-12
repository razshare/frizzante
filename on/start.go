package on

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	flag "github.com/spf13/pflag"
	"os"
)

func Start(efs embed.FS) {
	if !state.Parsed {
		flag.Parse()
		state.Parsed = true
	}

	if *state.Help {
		Help()
		os.Exit(0)
	}

	if *state.Version {
		Version(efs)
		os.Exit(0)
	}

	if *state.CreateProject != "" {
		CreateProject(*state.CreateProject)
		os.Exit(0)
	}

	if *state.Generate != "" {
		Generate(efs, *state.Generate)
		os.Exit(0)
	}

	if *state.Test {
		Test()
		os.Exit(0)
	}

	if *state.Package {
		Package()
		os.Exit(0)
	}

	if *state.PackageWatch {
		PackageWatch()
		os.Exit(0)
	}

	if *state.Check {
		Check()
		os.Exit(0)
	}

	if *state.Update {
		Update()
		os.Exit(0)
	}

	if *state.Install {
		Install()
		os.Exit(0)
	}

	if *state.Format {
		Format()
		os.Exit(0)
	}

	if *state.Touch {
		Touch()
		os.Exit(0)
	}

	if *state.Clean {
		Clean()
		os.Exit(0)
	}

	if *state.Dev {
		Dev()
		os.Exit(0)
	}

	if *state.Build {
		Build()
		os.Exit(0)
	}

	if *state.Configure {
		Configure(efs)
		os.Exit(0)
	}

	if *state.Welcome {
		Welcome()
		os.Exit(0)
	}

	Menu(efs)
}
