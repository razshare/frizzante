package apps

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/extensions"
	flag "github.com/spf13/pflag"
)

func New() *App {
	add := flag.StringP("add", "a", "", "adds packages")
	hlp := flag.BoolP("help", "h", false, "shows this help document")
	lck := flag.BoolP("lock-packages", "", false, "locks packages to the current version")
	ver := flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
	res := flag.BoolP("reset", "", false, "deletes frizzante global directory")
	crt := flag.StringP("create-project", "c", "", "creates a frizzante project")
	gen := flag.StringP("generate", "g", "", "generates code and resources")
	mig := flag.StringP("migrate", "m", "", "migrates database schema, requires a database connection string or file")
	tst := flag.BoolP("test", "t", false, "runs tests")
	pkg := flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
	pkw := flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
	chk := flag.BoolP("check", "", false, "checks source code for errors")
	upd := flag.BoolP("update", "u", false, "updates go and js packages")
	ins := flag.BoolP("install", "i", false, "installs go and js packages")
	fmt := flag.BoolP("format", "f", false, "formats source code")
	tch := flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
	cln := flag.BoolP("clean-project", "", false, "deletes .gen, .vite, app/{dist,node_modules}")
	dev := flag.BoolP("dev", "d", false, "starts dev mode")
	bld := flag.BoolP("build", "b", false, "builds project")
	cnf := flag.BoolP("configure", "", false, "configures project by installing required binaries and packages")
	yes := flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
	wel := flag.BoolP("welcome", "", false, "shows a welcome message")
	clr := flag.BoolP("clear", "", false, "clears screen")
	asm := flag.BoolP("assembly-explorer", "", false, "shows the assembly explorer")
	_go := flag.StringP("go", "", "go"+extensions.Find(), "sets the go binary")
	air := flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extensions.Find()), "sets the air binary")
	bun := flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extensions.Find()), "sets the bun binary")
	sqc := flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()), "sets the sqlc binary")
	sqy := flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file")
	tgs := flag.StringP("tags", "", "", "sets build tags")
	dbs := flag.StringP("database", "", "", "specifies the database connection string or file, used when migrating database schema")
	snp := flag.BoolP("snapshot", "", false, "snapshots the server state and generates static web assets")

	return &App{
		Add:              add,
		Help:             hlp,
		Version:          ver,
		Reset:            res,
		CreateProject:    crt,
		Generate:         gen,
		Migrate:          mig,
		Test:             tst,
		Package:          pkg,
		PackageWatch:     pkw,
		Check:            chk,
		Update:           upd,
		Install:          ins,
		Format:           fmt,
		Touch:            tch,
		CleanProject:     cln,
		Dev:              dev,
		Build:            bld,
		Configure:        cnf,
		Yes:              yes,
		Go:               _go,
		Air:              air,
		Bun:              bun,
		Sqlc:             sqc,
		SqlcYaml:         sqy,
		Welcome:          wel,
		Tags:             tgs,
		AssemblyExplorer: asm,
		Clear:            clr,
		Database:         dbs,
		LockPackages:     lck,
		Snapshot:         snp,
	}
}
