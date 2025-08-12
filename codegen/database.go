package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Database(efs embed.FS) {
	t := strings.ToLower(singleselect.Send(
		[]string{"Sqlite"},
		"What type of database would you like to setup?",
	))

	to := filepath.Join("lib", "database")

	err := Generate(efs, []Generation{
		{
			From: "template/lib/database/" + t,
			To:   to,
			Overwrite: func(n string) bool {
				yes := confirm.Sendf(true, "file `%s` already exists. Overwrite?", n)
				if yes {
					messages.Infof("overwriting file `%s`", n)
				} else {
					messages.Infof("skipping file `%s`", n)
				}
				return yes
			},
		},
	})

	if err != nil {
		messages.Fatal(err)
		return
	}

	if t == "sqlite" {
		s := spinner.New("adding github.com/mattn/go-sqlite3")
		err = spinner.Start(s)
		if err != nil {
			messages.Fatal(err)
			return
		}

		install := exec.Command(path.Go("."), "get", "github.com/mattn/go-sqlite3")
		install.Env = append(os.Environ())
		install.Stderr = os.Stderr
		install.Stdout = os.Stdout
		install.Stdin = os.Stdin
		err = install.Run()
		if err != nil {
			spinner.Stop(s)
			messages.Fatal(err)
		}
		spinner.Stop(s)

		s = spinner.New("updating go dependencies")
		err = spinner.Start(s)
		if err != nil {
			messages.Fatal(err)
			return
		}
		get := exec.Command(path.Go("."), "get", "-u", "./...")
		get.Env = append(os.Environ())
		get.Stderr = os.Stderr
		get.Stdout = os.Stdout
		get.Stdin = os.Stdin
		err = get.Run()
		if err != nil {
			spinner.Stop(s)
			messages.Fatal(err)
		}
		spinner.Stop(s)
		messages.Success("your sqlite database is ready")

		if strings.Contains(strings.ToLower(*state.Generate), "queries") &&
			confirm.Send(true, "would you like to also generate your queries?") {
			Queries(efs)
		}
	}
}
