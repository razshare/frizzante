package cli

import (
	flag "github.com/spf13/pflag"
	"path/filepath"
)

var App = flag.StringP("app", "", "app", "sets the app directory")
var Help = flag.BoolP("help", "h", false, "shows this help document")
var Version = flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
var CreateProject = flag.StringP("create-project", "c", "", "creates a frizzante project")
var Generate = flag.StringP("generate", "g", "", "generates code and binaries")
var Test = flag.BoolP("test", "t", false, "runs tests")
var Package = flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
var PackageWatch = flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
var Check = flag.BoolP("check", "", false, "checks source code for errors")
var Update = flag.BoolP("update", "u", false, "updates dependencies")
var Install = flag.BoolP("install", "i", false, "installs dependencies")
var Format = flag.BoolP("format", "f", false, "formats source code")
var Touch = flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
var Clean = flag.BoolP("clean", "", false, "cleans project")
var Dev = flag.BoolP("dev", "d", false, "starts dev mode")
var Build = flag.BoolP("build", "b", false, "builds project")
var Configure = flag.BoolP("configure", "", false, "configures project by installing necessary binaries under \"./.gen\"")
var Platform = flag.StringP("platform", "", "", "sets the platform, accepts \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\", \"darwin/amd64\", \"windows/arm64\", \"windows/amd64\"")
var Yes = flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
var Go = flag.StringP("go", "", "go", "sets the go binary")
var Air = flag.StringP("air", "", filepath.Join(".gen", "air", "air"), "sets the air binary")
var Bun = flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"), "sets the bun binary")
var Sqlc = flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"), "sets the sqlc binary")
var Welcome = flag.BoolP("welcome", "", false, "shows a welcome message")
var Clear = flag.BoolP("clear", "", false, "clears screen on start")
