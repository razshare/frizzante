package apps

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/extension"
	flag "github.com/spf13/pflag"
)

func New() *App {
	add := flag.StringP("add", "a", "", "adds packages")
	app := flag.StringP("app", "", "app", "sets the app directory")
	help := flag.BoolP("help", "h", false, "shows this help document")
	ver := flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
	res := flag.BoolP("reset", "", false, "deletes frizzante global directory")
	crt := flag.StringP("create-project", "c", "", "creates a frizzante project")
	gen := flag.StringP("generate", "g", "", "generates code and resources")
	tst := flag.BoolP("test", "t", false, "runs tests")
	pkg := flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
	pkgw := flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
	chk := flag.BoolP("check", "", false, "checks source code for errors")
	upd := flag.BoolP("update", "u", false, "updates go and js packages")
	ins := flag.BoolP("install", "i", false, "installs go and js packages")
	fmt := flag.BoolP("format", "f", false, "formats source code")
	tch := flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
	cln := flag.BoolP("clean-project", "", false, "deletes .gen, .vite, app/{dist,node_modules}")
	dev := flag.BoolP("dev", "d", false, "starts dev mode")
	bld := flag.BoolP("build", "b", false, "builds project")
	cnf := flag.BoolP("configure", "", false, "configures project by installing required binaries and packages")
	plt := flag.StringP("platform", "", "", "sets the platform, accepts \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\", \"darwin/amd64\", \"windows/arm64\", \"windows/amd64\"")
	yes := flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
	_go := flag.StringP("go", "", "go"+extension.Find(), "sets the go binary")
	air := flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extension.Find()), "sets the air binary")
	bun := flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extension.Find()), "sets the bun binary")
	sqc := flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extension.Find()), "sets the sqlc binary")
	sqy := flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file")
	wel := flag.BoolP("welcome", "", false, "shows a welcome message")
	clr := flag.BoolP("clear", "", false, "clears screen")
	tags := flag.StringP("tags", "", "", "sets build tags")
	asme := flag.BoolP("assembly-explorer", "", false, "shows the assembly explorer")
	mig := flag.BoolP("migrate", "", false, "migrates database schema")
	mii := flag.IntP("migrate-index", "", 0, "index to which to migrate database schema")

	return &App{
		Add:              add,
		App:              app,
		Help:             help,
		Version:          ver,
		Reset:            res,
		CreateProject:    crt,
		Generate:         gen,
		Test:             tst,
		Package:          pkg,
		PackageWatch:     pkgw,
		Check:            chk,
		Update:           upd,
		Install:          ins,
		Format:           fmt,
		Touch:            tch,
		CleanProject:     cln,
		Dev:              dev,
		Build:            bld,
		Configure:        cnf,
		Platform:         plt,
		Yes:              yes,
		Go:               _go,
		Air:              air,
		Bun:              bun,
		Sqlc:             sqc,
		SqlcYaml:         sqy,
		Welcome:          wel,
		Tags:             tags,
		AssemblyExplorer: asme,
		Clear:            clr,
		Migrate:          mig,
		MigrateIndex:     mii,
	}
}
