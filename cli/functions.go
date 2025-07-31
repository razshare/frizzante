package cli

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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
var FlagHooks = flag.BoolP("hooks", "", false, "adds git hooks")
var FlagConfigure = flag.BoolP("configure", "", false, "configures project by installing necessary binaries under \"./.gen\"")
var FlagPlatform = flag.StringP("platform", "", "", "sets the platform, accepts either \"linux/amd64\", \"linux/arm64\", \"darwin/arm64\" or \"darwin/amd64\"")
var FlagYes = flag.BoolP("yes", "y", false, "confirms all binary prompts silently")
var FlagGo = flag.StringP("go", "", "go", "sets the go binary")
var FlagAir = flag.StringP("air", "", filepath.Join(".gen", "air", "air"), "sets the air binary")
var FlagBun = flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"), "sets the bun binary")
var FlagSqlc = flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"), "sets the sqlc binary")
var FlagSqlcGenerate = flag.BoolP("sqlc-generate", "", false, "generates sqlc queries")
var FlagWelcome = flag.BoolP("welcome", "", false, "shows a welcome message")

func OnStart(cli *Cli) {
	if !cli.Parsed {
		flag.Parse()
		cli.Parsed = true
	}

	if *FlagHelp {
		OnHelp(cli)
		os.Exit(0)
	}

	if *FlagVersion {
		OnVersion(cli)
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		OnCreateProject(cli, *FlagCreateProject)
		os.Exit(0)
	}

	if *FlagAdd != "" {
		OnAddFeature(cli, *FlagAdd)
		os.Exit(0)
	}

	if *FlagTest {
		OnTest(cli)
		os.Exit(0)
	}

	if *FlagPackage {
		OnPackage(cli)
		os.Exit(0)
	}

	if *FlagPackageWatch {
		OnPackageWatch(cli)
		os.Exit(0)
	}

	if *FlagCheck {
		OnCheck(cli)
		os.Exit(0)
	}

	if *FlagUpdate {
		OnUpdate(cli)
		os.Exit(0)
	}

	if *FlagInstall {
		OnInstall(cli)
		os.Exit(0)
	}

	if *FlagFormat {
		OnFormat(cli)
		os.Exit(0)
	}

	if *FlagTouch {
		OnTouch(cli)
		os.Exit(0)
	}

	if *FlagClean {
		OnClean(cli)
		os.Exit(0)
	}

	if *FlagDev {
		OnDev(cli)
		os.Exit(0)
	}

	if *FlagBuild {
		OnBuild(cli)
		os.Exit(0)
	}

	if *FlagHooks {
		OnHooks(cli)
		os.Exit(0)
	}

	if *FlagConfigure {
		OnConfigure(cli)
		os.Exit(0)
	}

	if *FlagSqlcGenerate {
		OnSqlcGenerate(cli)
		os.Exit(0)
	}

	if *FlagWelcome {
		OnWelcome(cli)
		os.Exit(0)
	}

	OnMenu(cli)
}

func OnMenu(cli *Cli) {
	err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("Frizzante", pterm.FgCyan.ToStyle())).Render()
	if err != nil {
		Fatal(cli, err)
	}

	options := []string{
		"Help",
		"Update",
		"Install",
		"Version",
		"Create Project",
		"Add",
		"Add?",
		"Test",
		"Package",
		"Package Watch",
		"Check",
		"Format",
		"Touch",
		"Clean",
		"Dev",
		"Build",
		"Hooks",
		"Configure",
		"Sqlc Generate",
	}

	result, showError := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Pick an option")
	if showError != nil {
		Fatal(cli, showError)
	}

	if result == "Help" {
		*FlagHelp = true
		OnStart(cli)
		return
	}

	if result == "Version" {
		*FlagVersion = true
		OnStart(cli)
		return
	}

	if result == "Create Project" {
		projectName, projectNameError := pterm.DefaultInteractiveTextInput.Show("Give the project name")
		if projectNameError != nil {
			Fatal(cli, projectNameError)
		}
		*FlagCreateProject = projectName
		OnStart(cli)
		return
	}

	if result == "Add" {
		*FlagAdd = ":pick"
		OnStart(cli)
		return
	}

	if result == "Add?" {
		*FlagAdd = "?"
		OnStart(cli)
		return
	}

	if result == "Test" {
		*FlagTest = true
		OnStart(cli)
		return
	}

	if result == "Package" {
		*FlagPackage = true
		OnStart(cli)
		return
	}

	if result == "Package Watch" {
		*FlagPackageWatch = true
		OnStart(cli)
		return
	}

	if result == "Check" {
		*FlagCheck = true
		OnStart(cli)
		return
	}

	if result == "Update" {
		*FlagUpdate = true
		OnStart(cli)
		return
	}

	if result == "Install" {
		*FlagInstall = true
		OnStart(cli)
		return
	}

	if result == "Format" {
		*FlagFormat = true
		OnStart(cli)
		return
	}

	if result == "Touch" {
		*FlagTouch = true
		OnStart(cli)
		return
	}

	if result == "Clean" {
		*FlagClean = true
		OnStart(cli)
		return
	}

	if result == "Dev" {
		*FlagDev = true
		OnStart(cli)
		return
	}

	if result == "Build" {
		*FlagBuild = true
		OnStart(cli)
		return
	}

	if result == "Hooks" {
		*FlagHooks = true
		OnStart(cli)
		return
	}

	if result == "Configure" {
		*FlagConfigure = true
		OnStart(cli)
		return
	}

	if result == "Sqlc Generate" {
		*FlagSqlcGenerate = true
		OnStart(cli)
		return
	}
}

func OnHelp(cli *Cli) {
	flag.Usage()
}

func OnVersion(cli *Cli) {
	var version string

	versionData, versionError := cli.Efs.ReadFile("version")
	if versionError != nil {
		Fatal(cli, versionError)
	}

	version = string(versionData)

	println(version)
}

func OnCreateProject(cli *Cli, project string) {
	downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", project+".zip")
	if downloadError != nil {
		Fatal(cli, downloadError)
	}

	unzipError := files.UnzipFile(project+".zip", project+".tmp")
	if unzipError != nil {
		Fatal(cli, unzipError)
	}

	removeError := os.Remove(project + ".zip")
	if removeError != nil {
		Fatal(cli, removeError)
	}

	renameError := os.Rename(filepath.Join(project+".tmp", "frizzante-starter-main"), project)
	if renameError != nil {
		Fatal(cli, renameError)
	}

	removeAllError := os.RemoveAll(filepath.Join(project + ".tmp"))
	if removeAllError != nil {
		Fatal(cli, removeAllError)
	}
}

func CopyFeatureDirectories(cli *Cli, instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		name := instruction.FeatureName
		from := instruction.OriginDirectory
		to := instruction.DestinationDirectory

		if files.IsDirectory(to) {
			if !Confirmf(cli, "It looks like feature `%s` already exists in this project, would you like to overwrite it?", name) {
				Infof(cli, "skipping `%s`", name)
				return
			}

			removeAllError := os.RemoveAll(to)
			if removeAllError != nil {
				Fatal(cli, removeAllError)
			}
		}

		unixFrom := strings.ReplaceAll(from, "\\", "/")
		unixFromFileNames, readDirError := embeds.ReadDir(cli.Efs, unixFrom)
		if readDirError != nil {
			Fatal(cli, readDirError)
		}

		for _, unixFileName := range unixFromFileNames {
			fileName := to + strings.ReplaceAll(strings.TrimPrefix(unixFileName, unixFrom), "/", string(filepath.Separator))
			directoryName := filepath.Dir(fileName)

			if !files.IsDirectory(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if mkdirError != nil {
					Fatal(cli, mkdirError)
				}
			}

			file, openError := os.Create(fileName)
			if openError != nil {
				Fatal(cli, openError)
			}

			esfFile, esfOpenError := cli.Efs.Open(unixFileName)
			if esfOpenError != nil {
				Fatal(cli, esfOpenError)
			}

			_, copyError := io.Copy(file, esfFile)
			if copyError != nil {
				Fatal(cli, copyError)
			}
		}

		Successf(cli, "adding `%s`", to)
	}
}

func AddFeatureByName(cli *Cli, feature string) {
	if strings.ToLower(feature) == "core" {
		core := filepath.Join("app", "frizzante", "core")
		CopyFeatureDirectories(cli, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
		return
	}

	if strings.ToLower(feature) == "forms" {
		core := filepath.Join("app", "frizzante", "core")
		forms := filepath.Join("app", "frizzante", "forms")

		if !files.IsDirectory(core) {
			if Confirmf(cli, "It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				CopyFeatureDirectories(cli, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		CopyFeatureDirectories(cli, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: forms, DestinationDirectory: forms}})
		return
	}

	if strings.ToLower(feature) == "links" {
		core := filepath.Join("app", "frizzante", "core")
		links := filepath.Join("app", "frizzante", "links")

		if !files.IsDirectory(core) {
			if Confirmf(cli, "It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				CopyFeatureDirectories(cli, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		CopyFeatureDirectories(cli, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: links, DestinationDirectory: links}})
		return
	}

	if strings.ToLower(feature) == "bun" {
		directoryName := filepath.Join(".gen", "bun")

		platform := Platform(cli)

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
		} else if platform == PlatformLinuxArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
		}

		Install(cli, "bun", url, directoryName)

		var fileName string

		if platform == PlatformDarwinArm64 {
			fileName = filepath.Join(directoryName, "bun-darwin-aarch64", "bun")
		} else if platform == PlatformDarwinAmd64 {
			fileName = filepath.Join(directoryName, "bun-darwin-x64", "bun")
		} else if platform == PlatformLinuxArm64 {
			fileName = filepath.Join(directoryName, "bun-linux-aarch64", "bun")
		} else if platform == PlatformLinuxAmd64 {
			fileName = filepath.Join(directoryName, "bun-linux-x64", "bun")
		}

		if files.IsFile(fileName) {
			renameError := os.Rename(fileName, Bun(cli, "."))
			if renameError != nil {
				Fatal(cli, renameError)
			}
		}

		installDirectory := filepath.Dir(fileName)
		if files.IsDirectory(directoryName) {
			removeError := os.RemoveAll(installDirectory)
			if removeError != nil {
				Fatal(cli, removeError)
			}
		}

		return
	}

	if strings.ToLower(feature) == "air" {
		directoryName := filepath.Join(".gen", "air")

		platform := Platform(cli)

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
		} else if platform == PlatformLinuxArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
		}

		Install(cli, "air", url, directoryName)

		return
	}

	if strings.ToLower(feature) == "sqlc" {
		directoryName := filepath.Join(".gen", "sqlc")

		platform := Platform(cli)

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if platform == PlatformLinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		}

		Install(cli, "sqlc", url, directoryName)

		writeSchemaSql := true
		writeQueriesSql := true
		writeSqlcYaml := true

		if files.IsFile("schema.sql") {
			writeSchemaSql = Confirm(cli, "File `schema.sql` already exists, would you like to overwrite it?")
		}

		if writeSchemaSql {
			writeError := os.WriteFile("schema.sql", make([]byte, 0), os.ModePerm)
			if writeError != nil {
				Fatal(cli, writeError)
			}
			Success(cli, "schema.sql created")
		}

		if files.IsFile("queries.sql") {
			writeQueriesSql = Confirm(cli, "File `queries.sql` already exists, would you like to overwrite it?")
		}

		if writeQueriesSql {
			writeError := os.WriteFile("queries.sql", make([]byte, 0), os.ModePerm)
			if writeError != nil {
				Fatal(cli, writeError)
			}
			Success(cli, "queries.sql created")
		}

		if files.IsFile("sqlc.yaml") {
			writeSqlcYaml = Confirm(cli, "File `sqlc.yaml` already exists, would you like to overwrite it?")
		}

		if writeSqlcYaml {
			data, readError := cli.Efs.ReadFile("sqlc.yaml")
			if readError != nil {
				Fatal(cli, readError)
			}

			readError = os.WriteFile("sqlc.yaml", data, os.ModePerm)
			if readError != nil {
				Fatal(cli, readError)
			}
			Success(cli, "sqlc.yaml created")
		}

		return
	}

	Fatalf(cli, "unknown feature `%s`", feature)
}

func Cwd(cli *Cli) string {
	wd, wdDir := os.Getwd()
	if wdDir != nil {
		Fatal(cli, wdDir)
	}
	return wd
}

func OnAddFeature(cli *Cli, features string) {
	if features == "?" {
		ShowFeaturesInfo(cli)
		return
	}

	if features == ":pick" {
		selectedFeatures, showError := pterm.
			DefaultInteractiveMultiselect.
			WithKeySelect(keys.Space).
			WithKeyConfirm(keys.Enter).
			WithFilter(true).
			WithOptions([]string{
				"Core",
				"Forms",
				"Links",
				"Air",
				"Bun",
				"Sqlc",
			}).
			WithFilter(false).
			Show("Pick a feature to add")

		if showError != nil {
			Fatal(cli, showError)
		}

		for _, selectedFeature := range selectedFeatures {
			AddFeatureByName(cli, selectedFeature)
		}
		return
	}

	splitFeatures := strings.Split(features, ",")

	for _, feature := range splitFeatures {
		AddFeatureByName(cli, feature)
	}
	return
}

func OnTest(cli *Cli) {
	OnPackage(cli)

	test := exec.Command(Go(cli, "."), "test")
	test.Env = append(os.Environ(), "CGO_ENABLED=1")
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	runError := test.Run()
	if runError != nil {
		Fatal(cli, runError)
	}
}

func OnHooks(cli *Cli) {
	fileName := ".git/hooks/pre-commit"
	directoryName := filepath.Dir(fileName)
	if !files.IsDirectory(directoryName) {
		Fatalf(cli, "directory `%s` not found", directoryName)
		return
	}

	if files.IsFile(fileName) {
		if !Confirm(cli, "This git repository already defines a pre-commit script, would you like to overwrite it?") {
			pterm.Info.Println("pre-commit hook skipped")
			return
		}
		removeError := os.Remove(fileName)
		if removeError != nil {
			Fatal(cli, removeError)
		}
		Success(cli, "pre-commit script overwritten")
	}

	err := os.WriteFile(fileName, []byte("make test"), os.ModePerm)
	if err != nil {
		Fatal(cli, err)
	}

	Success(cli, "hooks added")
}

func OnTouch(cli *Cli) {
	touch := func(fileName string) {
		directoryName := filepath.Dir(fileName)

		if !files.IsDirectory(directoryName) {
			mkdirAllError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirAllError != nil {
				Fatal(cli, mkdirAllError)
			}
		}

		file, openError := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if openError != nil {
			Fatal(cli, openError)
		}

		closeError := file.Close()
		if closeError != nil {
			Fatal(cli, closeError)
		}
	}

	mkdirError := os.MkdirAll("app/dist", os.ModePerm)
	if mkdirError != nil {
		Fatal(cli, mkdirError)
	}

	touch(filepath.Join("app", "dist", "server.js"))
	touch(filepath.Join("app", "dist", "client", "index.html"))
}

func OnClean(cli *Cli) {
	clean := exec.Command(Go(cli, "."), "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	runError := clean.Run()
	if runError != nil {
		Fatal(cli, runError)
	}

	removeError := os.RemoveAll(filepath.Join("app", "dist"))
	if removeError != nil {
		Fatal(cli, removeError)
	}

	removeError = os.RemoveAll(filepath.Join("app", "node_modules"))
	if removeError != nil {
		Fatal(cli, removeError)
	}

	removeError = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if removeError != nil {
		Fatal(cli, removeError)
	}

	removeError = os.RemoveAll(".vite")
	if removeError != nil {
		Fatal(cli, removeError)
	}

	OnTouch(cli)

	Success(cli, "project cleaned")
}

func OnFormat(cli *Cli) {
	OnTouch(cli)

	gofmt := exec.Command(Go(cli, "."), "fmt")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	gofmtError := gofmt.Run()
	if gofmtError != nil {
		Fatal(cli, gofmtError)
	}

	prettier := exec.Command(Bun(cli, "app"), "x", "prettier", "--write", ".")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		Fatal(cli, prettierError)
	}

	Success(cli, "project formatted")
}

func OnUpdate(cli *Cli) {
	OnTouch(cli)

	get := exec.Command(Go(cli, "."), "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	getError := get.Run()
	if getError != nil {
		Fatal(cli, getError)
	}

	prettier := exec.Command(Bun(cli, "app"), "update")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		Fatal(cli, prettierError)
	}

	Success(cli, "project dependencies updated")
}

func OnInstall(cli *Cli) {
	OnTouch(cli)

	tidy := exec.Command(Go(cli, "."), "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	tidyError := tidy.Run()
	if tidyError != nil {
		Fatal(cli, tidyError)
	}

	install := exec.Command(Bun(cli, "app"), "install")
	install.Dir = "app"
	install.Env = append(os.Environ())
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin
	installError := install.Run()
	if installError != nil {
		Fatal(cli, installError)
	}

	Success(cli, "project dependencies installed")
}

func OnPackage(cli *Cli) {
	OnTouch(cli)

	server := exec.Command(Bun(cli, "app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Run()
	if serverError != nil {
		Fatal(cli, serverError)
	}

	client := exec.Command(Bun(cli, "app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Run()
	if clientError != nil {
		Fatal(cli, clientError)
	}

	esbuild := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esbuild.Dir = "app"
	esbuild.Env = append(os.Environ())
	esbuild.Stderr = os.Stderr
	esbuild.Stdout = os.Stdout
	esbuild.Stdin = os.Stdin
	esbuildError := esbuild.Run()
	if esbuildError != nil {
		Fatal(cli, esbuildError)
	}

	Success(cli, "project app package generated in app/dist")
}

func OnPackageWatch(cli *Cli) {
	OnTouch(cli)

	server := exec.Command(Bun(cli, "app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Start()
	if serverError != nil {
		Fatalf(cli, "vite server watcher failed to launch\n%s", serverError)
	}
	Success(cli, "vite server watcher launched")

	client := exec.Command(Bun(cli, "app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Start()
	if clientError != nil {
		Fatalf(cli, "vite client watcher failed to launch\n%s", clientError)
	}
	Success(cli, "vite client watcher launched")

	clientWaitError := client.Wait()
	if clientWaitError != nil {
		Fatal(cli, clientWaitError)
	}

	serverWaitError := server.Wait()
	if serverWaitError != nil {
		Fatal(cli, serverWaitError)
	}
}

func OnDev(cli *Cli) {
	OnTouch(cli)

	mkdirError := os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if mkdirError != nil {
		Fatal(cli, mkdirError)
	}

	air := exec.Command(Air(cli, "."))
	air.Env = append(os.Environ(), "DEV=1", "CGO_ENABLED=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	airError := air.Start()
	if airError != nil {
		Fatalf(cli, "air watcher faield to launch\n%s", airError)
	}
	Success(cli, "air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { OnPackageWatch(cli) }()

	group.Wait()
	tidyWaitError := air.Wait()
	if tidyWaitError != nil {
		Fatal(cli, tidyWaitError)
	}
}

func OnBuild(cli *Cli) {
	OnPackage(cli)

	build := exec.Command(Go(cli, "."), "build", "-o=.gen/bin/app", ".")
	build.Env = append(os.Environ(), "CGO_ENABLED=1")

	if strings.ToLower(*FlagPlatform) == "linux/amd64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=amd64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=linux", "GOARCH=arm64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=amd64")
	} else if strings.ToLower(*FlagPlatform) == "linux/arm64" {
		build.Env = append(build.Env, "GOOS=darwin", "GOARCH=arm64")
	}

	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	buildError := build.Run()
	if buildError != nil {
		Fatal(cli, buildError)
	}
	Success(cli, "project built into .gen/bin/app")
}

func OnCheck(cli *Cli) {
	OnTouch(cli)

	eslint := exec.Command(Bun(cli, "app"), "x", "eslint")
	eslint.Dir = "app"
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	eslintError := eslint.Run()
	if eslintError != nil {
		Fatal(cli, eslintError)
	}

	svelteCheck := exec.Command(Bun(cli, "app"), "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = "app"
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	svelteCheckError := svelteCheck.Run()
	if svelteCheckError != nil {
		Fatal(cli, svelteCheckError)
	}
}

func OnConfigure(cli *Cli) {
	OnAddFeature(cli, "bun,air")
	OnInstall(cli)
}

func OnSqlcGenerate(cli *Cli) {
	sqlcGenerate := exec.Command(Sqlc(cli, "."), "generate")
	sqlcGenerate.Env = append(os.Environ())
	sqlcGenerate.Stderr = os.Stderr
	sqlcGenerate.Stdout = os.Stdout
	sqlcGenerate.Stdin = os.Stdin
	svelteCheckError := sqlcGenerate.Run()
	if svelteCheckError != nil {
		Fatal(cli, svelteCheckError)
	}
}

func OnWelcome(cli *Cli) {
	usingDocker := os.Getenv("FRIZZANTE_USING_DOCKER")
	end := make(chan string, 0)
	err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("Frizzante", pterm.FgCyan.ToStyle())).Render()
	if err != nil {
		Fatal(cli, err)
	}

	if usingDocker != "" {
		// TODO: add a custom message for docker.
	}

	<-end
	Success(cli, "Bye!")
}

func Install(cli *Cli, name string, url string, destination string) {
	if files.IsDirectory(destination) {
		if !Confirmf(cli, "It looks like `%s` is already installed in `%s`, would you like to overwrite it?", name, destination) {
			Infof(cli, "skipping `%s`", name)
			return
		}

		removeError := os.RemoveAll(destination)
		if removeError != nil {
			Fatal(cli, removeError)
		}
	}

	spinner, spinnerError := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf("installing `%s` from `%s`...", name, url))
	if spinnerError != nil {
		Fatal(cli, spinnerError)
	}
	defer func() {
		stopError := spinner.Stop()
		if stopError != nil {
			Fatal(cli, stopError)
		}
	}()

	if !strings.HasSuffix(url, ".zip") {
		downloadError := files.DownloadFile(url, filepath.Join(destination, name))
		if downloadError != nil {
			Fatal(cli, downloadError)
		}

		Successf(cli, "%s installed in `%s`", name, destination)
		return
	}

	zipFileName := destination + ".zip"
	downloadError := files.DownloadFile(url, zipFileName)
	if downloadError != nil {
		Fatal(cli, downloadError)
	}
	defer func() {
		removeError := os.Remove(zipFileName)
		if removeError != nil {
			Fatal(cli, removeError)
		}
	}()

	unzipError := files.UnzipFile(zipFileName, destination)
	if unzipError != nil {
		Fatal(cli, unzipError)
	}

	Successf(cli, "%s installed in `%s`", name, destination)
}

func ShowFeaturesInfo(cli *Cli) {
	Info(cli, strings.Join([]string{
		"You can use -a or --add",
		"in order to add new features to the project.",
		"",
		"The value passed in must follow",
		"the syntax: `-a{feature},{feature}`",
		"where {feature} is the name of the feature.",
		"",
		"For example, `-acore,forms` will generate the core and forms",
		"features of frizzante respectively in `app/frizzante/core` and `app/frizzante/forms`.",
		"",
		"Feature names are not case-sensitive.",
		"",
		"You can also use -a:pick or --add :pick to pick feature interactively.",
	}, "\n"))

	pterm.Println()

	data := pterm.TableData{
		{"Feature Name", "Description"},
		{
			"Core",
			strings.Join([]string{
				"Adds the core of frizzante.",
				"A bundle of scripts and components that manage",
				"view rendering, view transitions, automatic state management,",
				"provides commonly used functions.",
				"",
				"Source code will be dropped in `app/frizzante/core`.",
			}, "\n"),
		},
		{
			"Forms",
			strings.Join([]string{
				"Adds a <Form> component which behaves like a <form> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
				"",
				"Source code will be dropped in `app/frizzante/forms`.",
				"",
				"Requires `Core`.",
			}, "\n"),
		},
		{
			"Links",
			strings.Join([]string{
				"Adds a <Link> component which behaves like an <a> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
				"",
				"Source code will be dropped in `app/frizzante/links`.",
				"",
				"Requires `Core`.",
			}, "\n"),
		},
		{
			"Bun",
			strings.Join([]string{
				"Adds bun to the project.",
				"",
				"Binaries will be dropped in `.gen/bun`.",
				"",
				"Bun is required for development mode.",
			}, "\n"),
		},
		{
			"Sqlc",
			strings.Join([]string{
				"Adds sqlc to the project.",
				"",
				"Binaries will be dropped in `.gen/sqlc`.",
			}, "\n"),
		},
	}

	// Create a table with a header and the defined data, then render it
	tableError := pterm.
		DefaultTable.
		WithHasHeader().
		WithData(data).
		WithBoxed(true).
		WithRowSeparator("─").
		WithHeaderRowSeparator("=").
		Render()
	if tableError != nil {
		Fatal(cli, tableError)
	}
}

type PlatformType uint

const PlatformLinuxAmd64 PlatformType = 0
const PlatformLinuxArm64 PlatformType = 1
const PlatformDarwinAmd64 PlatformType = 2
const PlatformDarwinArm64 PlatformType = 3

func Platform(cli *Cli) PlatformType {
	var platform string

	if *FlagPlatform != "" {
		platform = *FlagPlatform
	} else {
		var platformError error
		platform, platformError = pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Linux/amd64",
				"Linux/arm64",
				"Darwin/amd64",
				"Darwin/arm64",
			}).
			WithFilter(false).
			Show("Pick a platform")

		if platformError != nil {
			Fatal(cli, platformError)
		}
		*FlagPlatform = platform
	}

	if strings.ToLower(platform) == "linux/amd64" {
		return PlatformLinuxAmd64
	}

	if strings.ToLower(platform) == "linux/arm64" {
		return PlatformLinuxArm64
	}

	if strings.ToLower(platform) == "darwin/arm64" {
		return PlatformDarwinArm64
	}

	if strings.ToLower(platform) == "darwin/amd64" {
		return PlatformDarwinAmd64
	}

	Fatalf(cli, "unknown platform `%s`", platform)
	return PlatformLinuxAmd64 // Noop, cli.Fatalf will crash intentionally.
}

// Confirm shows a confirmation prompt.
//
// Returns true if the user confirms, otherwise false.
func Confirm(cli *Cli, text string) bool {
	if *FlagYes {
		return true
	}

	yes, showError := pterm.
		DefaultInteractiveConfirm.
		WithConfirmText("Y").
		WithDefaultText("n").
		WithDefaultValue(true).
		Show(text)

	if showError != nil {
		Fatal(cli, showError)
	}

	return yes
}

func Go(cli *Cli, basepath string) string {
	var goBinary string

	if *FlagGo != "" {
		goBinary = *FlagGo
	} else {
		goBinary = Go(cli, ".")
	}

	if strings.HasPrefix(goBinary, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		goBinary = strings.Replace(goBinary, "~", dirname, 1)
		return goBinary
	}

	if !strings.Contains(goBinary, string(filepath.Separator)) {
		return goBinary
	}

	path, pathError := filepath.Rel(basepath, goBinary)
	if pathError != nil {
		Fatal(cli, pathError)
	}

	return path
}

func Air(cli *Cli, basepath string) string {
	var air string

	if *FlagAir != "" {
		air = *FlagAir
	} else {
		air = filepath.Join(".gen", "air", "air")
	}

	if strings.HasPrefix(air, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		air = strings.Replace(air, "~", dirname, 1)
		return air
	}

	if !strings.Contains(air, string(filepath.Separator)) {
		return air
	}

	path, pathError := filepath.Rel(basepath, air)
	if pathError != nil {
		Fatal(cli, pathError)
	}

	return path
}

func Bun(cli *Cli, basepath string) string {
	var bun string

	if *FlagBun != "" {
		bun = *FlagBun
	} else {
		bun = filepath.Join(".gen", "bun", "bun")
	}

	if strings.HasPrefix(bun, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		bun = strings.Replace(bun, "~", dirname, 1)
		return bun
	}

	if !strings.Contains(bun, string(filepath.Separator)) {
		return bun
	}

	path, pathError := filepath.Rel(basepath, bun)
	if pathError != nil {
		Fatal(cli, pathError)
	}

	return path
}

func Sqlc(cli *Cli, basepath string) string {
	var sqlc string

	if *FlagSqlc != "" {
		sqlc = *FlagSqlc
	} else {
		sqlc = filepath.Join(".gen", "sqlc", "sqlc")
	}

	if strings.HasPrefix(sqlc, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		sqlc = strings.Replace(sqlc, "~", dirname, 1)
		return sqlc
	}

	if !strings.Contains(sqlc, string(filepath.Separator)) {
		return sqlc
	}

	path, pathError := filepath.Rel(basepath, sqlc)
	if pathError != nil {
		Fatal(cli, pathError)
	}

	return path
}

// Confirmf shows a confirmation prompt.
//
// Returns true if the user confirms, otherwise false.
func Confirmf(cli *Cli, template string, vars ...any) bool {
	return Confirm(cli, fmt.Sprintf(template, vars...))
}

// Fatalf shows a fatal message and terminates the application.
func Fatalf(cli *Cli, template string, vars ...any) {
	pterm.Fatal.Printfln(template, vars...)
}

// Warningf shows a warning message.
func Warningf(cli *Cli, template string, vars ...any) {
	pterm.Warning.Printfln(template, vars...)
}

// Infof shows an info message.
func Infof(cli *Cli, template string, vars ...any) {
	pterm.Info.Printfln(template, vars...)
}

// Successf shows a success message.
func Successf(cli *Cli, template string, vars ...any) {
	pterm.Success.Printfln(template, vars...)
}

// Fatal shows a fatal message and terminates the application.
func Fatal(cli *Cli, vars ...any) {
	pterm.Fatal.Println(vars...)
}

// Warning shows a warning message.
func Warning(cli *Cli, vars ...any) {
	pterm.Warning.Println(vars...)
}

// Info shows an info message.
func Info(cli *Cli, vars ...any) {
	pterm.Info.Println(vars...)
}

// Success shows a success message.
func Success(cli *Cli, vars ...any) {
	pterm.Success.Println(vars...)
}

// Section shows the name of a section using Markdown semantics.
func Section(cli *Cli, vars ...any) {
	pterm.DefaultSection.Println(vars...)
}
