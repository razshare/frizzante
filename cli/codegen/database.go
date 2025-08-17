package codegen

import (
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"github.com/razshare/frizzante/tui/text"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Database(gen string, clr bool, yes bool, gobin string, sqlcbin string) error {
	tchoice, err := singleselect.Send(
		[]string{"Sqlite"},
		"what type of database would you like to setup?",
	)

	if err != nil {
		return err
	}

	if clr {
		text.Clrscr()
	}

	t := strings.ToLower(tchoice)

	dst := filepath.Join("lib", "database")

	err = Copy([]CopyInstruction{
		{
			From: "template/lib/database/" + t,
			To:   dst,
			Overwrite: func(n string) (bool, error) {
				if yes {
					return true, nil
				}

				var overwrite bool
				overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", n)
				if err != nil {
					return false, err
				}

				if overwrite {
					messages.Infof("overwriting %s", n)
				} else {
					messages.Infof("skipping %s", n)
				}

				return overwrite, nil
			},
		},
	})

	if err != nil {
		return err
	}

	if t == "sqlite" {
		s := spinner.New("adding github.com/mattn/go-sqlite3")

		go spinner.Start(s)
		install := exec.Command(gobin, "get", "github.com/mattn/go-sqlite3")
		install.Env = append(os.Environ())
		install.Stderr = os.Stderr
		install.Stdout = os.Stdout
		install.Stdin = os.Stdin
		err = install.Run()
		spinner.Stop(s)

		if err != nil {
			return err
		}

		s = spinner.New("updating go dependencies")

		go spinner.Start(s)
		get := exec.Command(gobin, "get", "-u", "./...")
		get.Env = append(os.Environ())
		get.Stderr = os.Stderr
		get.Stdout = os.Stdout
		get.Stdin = os.Stdin
		err = get.Run()
		spinner.Stop(s)

		if err != nil {
			return err
		}

		messages.Success("sqlite database is ready")

		if strings.Contains(strings.ToLower(gen), "queries") {
			var queries bool
			queries, err = confirm.Send(true, "would you like to also generate your queries?")
			if err != nil {
				return err
			}
			if queries {
				err = Queries(clr, sqlcbin)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
