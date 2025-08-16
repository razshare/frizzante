package cli

import (
	flag "github.com/spf13/pflag"
	"path/filepath"
)

func New() *Cli {
	return &Cli{
		Flags: &Flags{
			App:           flag.StringP("app", "a", "app", "sets the location of the app directory"),
			Help:          flag.BoolP("help", "h", false, "shows this help document"),
			Version:       flag.BoolP("version", "v", false, "shows the frizzante version used by this binary"),
			CreateProject: flag.StringP("create-project", "c", "", "creates a frizzante project"),
			Generate:      flag.StringP("generate", "g", "", "generates code and binaries"),
			Test:          flag.BoolP("test", "t", false, "runs tests"),
			Package:       flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist"),
			PackageWatch:  flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist"),
			Check:         flag.BoolP("check", "", false, "checks source code for errors"),
			Update:        flag.BoolP("update", "u", false, "updates dependencies"),
			Install:       flag.BoolP("install", "i", false, "installs dependencies"),
			Format:        flag.BoolP("format", "f", false, "formats source code"),
			Touch:         flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)"),
			Clean:         flag.BoolP("clean", "", false, "cleans project"),
			Dev:           flag.BoolP("dev", "d", false, "starts dev mode"),
			Build:         flag.BoolP("build", "b", false, "builds project"),
			Configure:     flag.BoolP("configure", "", false, "configures project by installing necessary binaries under \"./.gen\""),
			Platform:      flag.StringP("platform", "", "", "sets the platform, accepts \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\", \"darwin/amd64\", \"windows/arm64\", \"windows/amd64\""),
			Yes:           flag.BoolP("yes", "y", false, "confirms all binary prompts silently"),
			Go:            flag.StringP("go", "", "go", "sets the go binary"),
			Air:           flag.StringP("air", "", filepath.Join(".gen", "air", "air"), "sets the air binary"),
			Bun:           flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"), "sets the bun binary"),
			Sqlc:          flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"), "sets the sqlc binary"),
			Welcome:       flag.BoolP("welcome", "", false, "shows a welcome message"),
		},
	}
}
