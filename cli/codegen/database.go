package codegen

import (
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"strings"
)

func Database(o DatabaseOptions) error {
	tchoice, err := singleselect.Send(
		[]string{"sqlite"},
		"what type of database would you like to setup?",
	)

	if err != nil {
		return err
	}

	t := strings.ToLower(tchoice)

	err = Copy([]CopyInstruction{
		{
			From: "template/lib/database/" + t,
			To:   o.Lib,
			Overwrite: func(n string) (bool, error) {
				if o.Auto {
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
		install := exec.Command(o.Go, "get", "github.com/mattn/go-sqlite3")
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
		get := exec.Command(o.Go, "get", "-u", "./...")
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

		if strings.Contains(strings.ToLower(o.Generate), "queries") {
			var queries bool
			queries, err = confirm.Send(true, "would you like to also generate your queries?")
			if err != nil {
				return err
			}
			if queries {
				err = Queries(QueriesOptions{
					Auto:     o.Auto,
					Sqlc:     o.Sqlc,
					Platform: o.Platform,
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
