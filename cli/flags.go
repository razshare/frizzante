package cli

import (
	flag "github.com/spf13/pflag"
	"path/filepath"
)

var FlagHelp = flag.BoolP("help", "h", false, "shows this help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
var FlagCreateProject = flag.StringP("create-project", "c", "", "creates a frizzante project")
var FlagAdd = flag.StringP("add", "a", "", "adds features, see  \"-a?\" or \"--add ?\" for more details")
var FlagTest = flag.BoolP("test", "t", false, "runs tests")
var FlagPackage = flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
var FlagPackageWatch = flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
var FlagCheck = flag.BoolP("check", "", false, "checks source code for errors")
var FlagUpdate = flag.BoolP("update", "u", false, "updates dependencies")
var FlagInstall = flag.BoolP("install", "i", false, "installs dependencies")
var FlagFormat = flag.BoolP("format", "f", false, "formats source code")
var FlagTouch = flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
var FlagClean = flag.BoolP("clean", "", false, "cleans project")
var FlagDev = flag.BoolP("dev", "d", false, "starts dev mode")
var FlagBuild = flag.BoolP("build", "b", false, "builds project")
var FlagConfigure = flag.BoolP("configure", "", false, "configures project by installing necessary binaries under \"./.gen\"")
var FlagPlatform = flag.StringP("platform", "", "", "sets the platform, accepts \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\", \"darwin/amd64\", \"windows/arm64\", \"windows/amd64\"")
var FlagYes = flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
var FlagGo = flag.StringP("go", "", "go", "sets the go binary")
var FlagAir = flag.StringP("air", "", filepath.Join(".gen", "air", "air"), "sets the air binary")
var FlagBun = flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"), "sets the bun binary")
var FlagSqlc = flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"), "sets the sqlc binary")
var FlagSqlcGenerate = flag.BoolP("sqlc-generate", "", false, "generates sqlc queries")
var FlagCreateSqliteDatabase = flag.BoolP("create-sqlite-database", "", false, "creates an empty sqlite database")
var FlagWelcome = flag.BoolP("welcome", "", false, "shows a welcome message")
