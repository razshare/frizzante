package on

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/text"
	flag "github.com/spf13/pflag"
)

func Start() error {
	text.Clrscr()

	if !state.Parsed {
		flag.Parse()
		state.Parsed = true
	}

	gobin, err := path.Go(*cli.Go)
	if err != nil {
		return err
	}

	bunbin, err := path.Bun(*cli.Bun, *cli.App)
	if err != nil {
		return err
	}

	airbin, err := path.Air(*cli.Air)
	if err != nil {
		return err
	}

	sqlcbin, err := path.Sqlc(*cli.Sqlc, ".")
	if err != nil {
		return err
	}

	if *cli.Help {
		return Help()
	}

	if *cli.Version {
		return Version()
	}

	if *cli.CreateProject != "" {
		return CreateProject(*cli.CreateProject)
	}

	if *cli.Generate != "" {
		return Generate(*cli.App, *cli.Generate, *cli.Clear, *cli.Yes, gobin, sqlcbin)
	}

	if *cli.Test {
		return Test(*cli.App, gobin, bunbin)
	}

	if *cli.Package {
		return Package(*cli.App, bunbin)
	}

	if *cli.PackageWatch {
		return PackageWatch(*cli.App, bunbin)
	}

	if *cli.Check {
		return Check(*cli.App, bunbin)
	}

	if *cli.Update {
		return Update(*cli.App, gobin, bunbin)
	}

	if *cli.Install {
		return Install(*cli.App, gobin, bunbin)
	}

	if *cli.Format {
		return Format(*cli.App, gobin, bunbin)
	}

	if *cli.Touch {
		return Touch(*cli.App)
	}

	if *cli.Clean {
		return Clean(*cli.App, gobin)
	}

	if *cli.Dev {
		return Dev(*cli.App, airbin, bunbin)
	}

	if *cli.Build {
		return Build(*cli.App, *cli.Platform, gobin, bunbin)
	}

	if *cli.Configure {
		return Configure(*cli.App, *cli.Clear, *cli.Yes, gobin, bunbin, sqlcbin)
	}

	if *cli.Welcome {
		return Welcome()
	}

	logo, err := cli.Efs.ReadFile("clilogo.txt")
	if err != nil {
		return err
	}
	fmt.Println(config.Styles.BigText.Render(string(logo)))

	return Menu()
}
