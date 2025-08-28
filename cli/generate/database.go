package generate

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"strings"
)

func Database(options DatabaseOptions) (err error) {
	if files.IsDirectory(options.Lib) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", options.Lib); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", options.Lib)
				return nil
			}
		}

		if err = os.RemoveAll(options.Lib); err != nil {
			return
		}
	}

	var choice string
	if choice, err = singleselect.Send([]search.Choice{{Id: "sqlite"}}, "what type of database would you like to setup?"); err != nil {
		return
	}

	choice = strings.ToLower(choice)

	if err = Copy(CopyOptions{From: "internal/template/lib/database/" + choice, To: options.Lib, Auto: options.Auto}); err != nil {
		return
	}

	if choice == "sqlite" {
		spin := spinner.New("adding github.com/mattn/go-sqlite3")

		go spinner.Start(spin)
		install := exec.Command(options.Go, "get", "github.com/mattn/go-sqlite3")
		install.Env = append(os.Environ())
		install.Stderr = os.Stderr
		install.Stdout = os.Stdout
		install.Stdin = os.Stdin
		err = install.Run()
		spinner.Stop(spin)

		if err != nil {
			return
		}

		spin = spinner.New("updating go dependencies")

		go spinner.Start(spin)
		get := exec.Command(options.Go, "get", "-u", "./...")
		get.Env = append(os.Environ())
		get.Stderr = os.Stderr
		get.Stdout = os.Stdout
		get.Stdin = os.Stdin
		err = get.Run()
		spinner.Stop(spin)

		if err != nil {
			return
		}

		messages.Success("sqlite database is ready")

		if strings.Contains(strings.ToLower(options.Generate), "queries") {
			var queries bool
			if queries, err = confirm.Send(true, "would you like to also generate your queries?"); err != nil {
				return
			}

			if queries {
				if err = Queries(QueriesOptions{Auto: options.Auto, Sqlc: options.Sqlc, Platform: options.Platform}); err != nil {
					return
				}
			}
		}
	}

	return
}
