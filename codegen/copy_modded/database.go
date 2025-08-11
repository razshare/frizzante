package copy_modded

import (
	"embed"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"os/exec"
	"path/filepath"
)

func Database(efs embed.FS) {
	err := embeds.Generate(efs, []codegen.Generation{
		{
			From: "template/lib/database",
			To:   filepath.Join("lib", "database"),
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
}
