package cli

import (
	"github.com/razshare/frizzante/cli/extension"
	"github.com/razshare/frizzante/tui/input"
	flag "github.com/spf13/pflag"
	"path/filepath"
)

func New() (c *Cli) {
	app := flag.StringP("app", "", "app", "sets the app directory")
	help := flag.BoolP("help", "h", false, "shows this help document")
	ver := flag.BoolP("version", "v", false, "shows the frizzante version used by this binary")
	res := flag.BoolP("reset", "", false, "deletes frizzante global directory")
	crt := flag.StringP("create-project", "c", "", "creates a frizzante project")
	gen := flag.StringP("generate", "g", "", "generates code and binaries")
	tst := flag.BoolP("test", "t", false, "runs tests")
	pkg := flag.BoolP("package", "p", false, "packages app, result will be dropped in app/dist")
	pkgw := flag.BoolP("package-watch", "", false, "watches and packages app, result will be dropped in app/dist")
	chk := flag.BoolP("check", "", false, "checks source code for errors")
	upd := flag.BoolP("update", "u", false, "updates dependencies")
	ins := flag.BoolP("install", "i", false, "installs dependencies")
	fmt := flag.BoolP("format", "f", false, "formats source code")
	tch := flag.BoolP("touch", "", false, "creates placeholders in app/dist (useful for go:embed)")
	cln := flag.BoolP("clean", "", false, "cleans project")
	dev := flag.BoolP("dev", "d", false, "starts dev mode")
	bld := flag.BoolP("build", "b", false, "builds project")
	cnf := flag.BoolP("configure", "", false, "configures project by installing necessary binaries under \"./.gen\"")
	plt := flag.StringP("platform", "", "", "sets the platform, accepts \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\", \"darwin/amd64\", \"windows/arm64\", \"windows/amd64\"")
	yes := flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
	_go := flag.StringP("go", "", "go"+extension.Find(), "sets the go binary")
	air := flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extension.Find()), "sets the air binary")
	bun := flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extension.Find()), "sets the bun binary")
	sqc := flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extension.Find()), "sets the sqlc binary")
	wel := flag.BoolP("welcome", "", false, "shows a welcome message")
	clr := flag.BoolP("clear", "", false, "clears screen on start")
	pkr := Menu{
		"configure": func() (func(), error) {
			*cnf = true
			return func() { *cnf = false }, nil
		},
		"create project": func() (func(), error) {
			var err error
			*crt, err = input.Send("give the project a name")
			if err != nil {
				return nil, err
			}
			return func() { *crt = "" }, nil
		},
		"dev": func() (func(), error) {
			*dev = true
			return func() { *dev = false }, nil
		},
		"build": func() (func(), error) {
			*bld = true
			return func() { *bld = false }, nil
		},
		"install": func() (func(), error) {
			*ins = true
			return func() { *ins = false }, nil
		},
		"update": func() (func(), error) {
			*upd = true
			return func() { *upd = false }, nil
		},
		"generate": func() (func(), error) {
			*gen = ":pick"
			return func() { *gen = "" }, nil
		},
		"package": func() (func(), error) {
			*pkg = true
			return func() { *pkg = false }, nil
		},
		"package (watch)": func() (func(), error) {
			*pkgw = true
			return func() { *pkgw = false }, nil
		},
		"test": func() (func(), error) {
			*tst = true
			return func() { *tst = false }, nil
		},
		"check": func() (func(), error) {
			*chk = true
			return func() { *chk = false }, nil
		},
		"format": func() (func(), error) {
			*fmt = true
			return func() { *fmt = false }, nil
		},
		"touch": func() (func(), error) {
			*tch = true
			return func() { *tch = false }, nil
		},
		"clean": func() (func(), error) {
			*cln = true
			return func() { *cln = false }, nil
		},
		"help": func() (func(), error) {
			*help = true
			return func() { *help = false }, nil
		},
		"version": func() (func(), error) {
			*ver = true
			return func() { *ver = false }, nil
		},
		"reset": func() (func(), error) {
			*res = true
			return func() { *res = false }, nil
		},
	}

	return &Cli{
		App:          app,
		Help:         help,
		Version:      ver,
		Reset:        res,
		Project:      crt,
		Generate:     gen,
		Test:         tst,
		Package:      pkg,
		PackageWatch: pkgw,
		Check:        chk,
		Update:       upd,
		Install:      ins,
		Format:       fmt,
		Touch:        tch,
		Clean:        cln,
		Dev:          dev,
		Build:        bld,
		Configure:    cnf,
		Platform:     plt,
		Yes:          yes,
		Go:           _go,
		Air:          air,
		Bun:          bun,
		Sqlc:         sqc,
		Welcome:      wel,
		Clear:        clr,
		Menu:         pkr,
	}
}
