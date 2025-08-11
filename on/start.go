package on

import (
	"embed"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/flags"
	flag "github.com/spf13/pflag"
	"os"
)

func Start(efs embed.FS) {
	if !cli.Parsed {
		flag.Parse()
		cli.Parsed = true
	}

	if *flags.Help {
		Help()
		os.Exit(0)
	}

	if *flags.Version {
		Version(efs)
		os.Exit(0)
	}

	if *flags.CreateProject != "" {
		CreateProject(*flags.CreateProject)
		os.Exit(0)
	}

	if *flags.Generate != "" {
		Generate(efs, *flags.Generate)
		os.Exit(0)
	}

	if *flags.Test {
		Test()
		os.Exit(0)
	}

	if *flags.Package {
		Package()
		os.Exit(0)
	}

	if *flags.PackageWatch {
		PackageWatch()
		os.Exit(0)
	}

	if *flags.Check {
		Check()
		os.Exit(0)
	}

	if *flags.Update {
		Update()
		os.Exit(0)
	}

	if *flags.Install {
		Install()
		os.Exit(0)
	}

	if *flags.Format {
		Format()
		os.Exit(0)
	}

	if *flags.Touch {
		Touch()
		os.Exit(0)
	}

	if *flags.Clean {
		Clean()
		os.Exit(0)
	}

	if *flags.Dev {
		Dev()
		os.Exit(0)
	}

	if *flags.Build {
		Build()
		os.Exit(0)
	}

	if *flags.Configure {
		Configure(efs)
		os.Exit(0)
	}

	if *flags.Welcome {
		Welcome()
		os.Exit(0)
	}

	Menu(efs)
}
