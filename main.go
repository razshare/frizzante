package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/additions/**
//go:embed internal/project/**
//go:embed internal/project/lib/core/ssr/.gitignore
//go:embed internal/project/.vscode/launch.json
//go:embed internal/project/app/.gitignore
//go:embed internal/project/app/.npmrc
//go:embed internal/project/app/.prettierrc
//go:embed internal/project/app/.prettierignore
var efs embed.FS

func main() {
	messages.Prefix = configs.Styles.Menu.PaddingRight(1).Render("│")
	app := apps.App{
		Efs:            efs,
		Go:             flag.StringP("go", "", fmt.Sprintf("go%s", extensions.Find()), "sets the go binary location"),
		Air:            flag.StringP("air", "", filepath.Join(".gen", "air", fmt.Sprintf("air%s", extensions.Find())), "sets the air binary location"),
		Bun:            flag.StringP("bun", "", filepath.Join(".gen", "bun", fmt.Sprintf("bun%s", extensions.Find())), "sets the bun binary location"),
		Sqlc:           flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", fmt.Sprintf("sqlc%s", extensions.Find())), "sets the sqlc binary location"),
		SqlcYaml:       flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file location"),
		Tags:           flag.StringP("tags", "", "", "sets tags; used for dev mode and building for production"),
		Database:       flag.StringP("database", "", "", "sets the database connection string; used for running migrations and snapshots"),
		DatabaseType:   flag.StringP("database-type", "", "", "sets the type of database to use; currently only sqlite is supported"),
		Create:         flag.StringP("create", "c", "", "creates a new project"),
		Add:            flag.StringP("add", "a", "", "adds packages"),
		Context:        flag.StringP("context", "", "js", "sets the context; used with --add; currently only \"js\" context is supported"),
		Migrate:        flag.StringP("migrate", "", "", "migrates database schema"),
		Generate:       flag.StringP("generate", "g", "", "generates code and resources"),
		Strict:         flag.BoolP("strict", "", false, "enables strict mode; program will stop if any flags are missing; used in ci/cd pipelines"),
		Dev:            flag.BoolP("dev", "d", false, "runs air and vite in parallel"),
		Snapshot:       flag.BoolP("snapshot", "s", false, "runs air (with tag snapshot_servers) and vite in parallel"),
		Configure:      flag.BoolP("configure", "", false, "generates binaries, installs packages and creates app/dist"),
		Install:        flag.BoolP("install", "i", false, "installs packages"),
		Update:         flag.BoolP("updates", "u", false, "updates packages"),
		Build:          flag.BoolP("build", "b", false, "builds project"),
		Asm:            flag.BoolP("asm", "", false, "starts the assembly explorer"),
		Package:        flag.BoolP("package", "", false, "packages the svelte application into app/dist"),
		PackageWatch:   flag.BoolP("package-watch", "", false, "packages the svelte application when source code changes"),
		Check:          flag.BoolP("check", "", false, "checks for code errors"),
		Format:         flag.BoolP("format", "f", false, "formats svelte and go code"),
		Touch:          flag.BoolP("touch", "", false, "adds placeholders in app/dist"),
		Clean:          flag.BoolP("clean", "", false, "deletes .gen, .vite, app/{dist,node_modules}"),
		Reset:          flag.BoolP("reset", "", false, "deletes global cache"),
		Clear:          flag.BoolP("clear", "", false, "clears terminal screen"),
		LockJsPackages: flag.BoolP("lock-packages", "", false, "locks packages to the current exact version"),
		Test:           flag.BoolP("test", "t", false, "runts tests"),
		Welcome:        flag.BoolP("welcome", "", false, "shows a welcome message and yields without killing the process"),
		Help:           flag.BoolP("help", "h", false, "shows the help menu"),
		Version:        flag.BoolP("version", "v", false, "shows the binary version"),
	}
	flag.Parse()
	if err := cli.Start(cli.StartOptions{App: app}); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
