package apps

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/extensions"
	flag "github.com/spf13/pflag"
)

func New() *App {
	add := flag.StringP("add", "a", "", "adds packages")
	help := flag.BoolP("help", "h", false, "shows this help document")
	lockPackages := flag.BoolP("lock-packages", "", false, "locks packages to the current version")
	version := flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
	reset := flag.BoolP("reset", "", false, "deletes frizzante global directory")
	createProject := flag.StringP("create-project", "c", "", "creates a frizzante project")
	generate := flag.StringP("generate", "g", "", "generates code and resources")
	migrate := flag.StringP("migrate", "m", "", "migrates database schema, requires a database connection string or file")
	test := flag.BoolP("test", "t", false, "runs tests")
	packageApp := flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
	packagesApPWatch := flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
	check := flag.BoolP("check", "", false, "checks source code for errors")
	update := flag.BoolP("update", "u", false, "updates go and js packages")
	install := flag.BoolP("install", "i", false, "installs go and js packages")
	format := flag.BoolP("format", "f", false, "formats source code")
	touch := flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
	cleanProject := flag.BoolP("clean-project", "", false, "deletes .gen, .vite, app/{dist,node_modules}")
	dev := flag.BoolP("dev", "d", false, "starts dev mode")
	build := flag.BoolP("build", "b", false, "builds project")
	configure := flag.BoolP("configure", "", false, "configures project by installing required binaries and packages")
	welcome := flag.BoolP("welcome", "", false, "shows a welcome message")
	clearScreen := flag.BoolP("clear", "", false, "clears screen")
	assemblyExplorer := flag.BoolP("assembly-explorer", "", false, "shows the assembly explorer")
	goBinary := flag.StringP("go", "", "go"+extensions.Find(), "sets the go binary")
	air := flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extensions.Find()), "sets the air binary")
	bun := flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extensions.Find()), "sets the bun binary")
	sqc := flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()), "sets the sqlc binary")
	sqlcYaml := flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file")
	tags := flag.StringP("tags", "", "", "sets build tags")
	database := flag.StringP("database", "", "", "specifies the database connection string or file, used when migrating database schema")
	databaseType := flag.StringP("database-type", "", "sqlc", "specifies the type of database (currently only sqlite is supported)")
	snapshot := flag.BoolP("snapshot", "", false, "snapshots the server state and generates static web assets")

	return &App{
		Add:              add,
		Help:             help,
		Version:          version,
		Reset:            reset,
		CreateProject:    createProject,
		Generate:         generate,
		Migrate:          migrate,
		Test:             test,
		Package:          packageApp,
		PackageWatch:     packagesApPWatch,
		Check:            check,
		Update:           update,
		Install:          install,
		Format:           format,
		Touch:            touch,
		CleanProject:     cleanProject,
		Dev:              dev,
		Build:            build,
		Configure:        configure,
		Go:               goBinary,
		Air:              air,
		Bun:              bun,
		Sqlc:             sqc,
		SqlcYaml:         sqlcYaml,
		Welcome:          welcome,
		Tags:             tags,
		AssemblyExplorer: assemblyExplorer,
		Clear:            clearScreen,
		Database:         database,
		DatabaseType:     databaseType,
		LockPackages:     lockPackages,
		Snapshot:         snapshot,
	}
}
