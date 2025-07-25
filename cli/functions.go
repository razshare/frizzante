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
var FlagAdd = flag.StringP("add", "a", "", fmt.Sprintf("adds features, see  \"-a?\" or \"--add ?\" for more details"))
var FlagTest = flag.BoolP("test", "t", false, fmt.Sprintf("runs tests"))
var FlagPackage = flag.BoolP("package", "p", false, fmt.Sprintf("packages app, result will be dropped in app/dist"))
var FlagPackageWatch = flag.BoolP("package-watch", "", false, fmt.Sprintf("watches and packages app, result will be dropped in app/dist"))
var FlagCheck = flag.BoolP("check", "", false, fmt.Sprintf("checks source code for errors"))
var FlagUpdate = flag.BoolP("update", "u", false, fmt.Sprintf("updates dependencies"))
var FlagInstall = flag.BoolP("install", "i", false, fmt.Sprintf("installs dependencies"))
var FlagFormat = flag.BoolP("format", "f", false, fmt.Sprintf("formats source code"))
var FlagTouch = flag.BoolP("touch", "", false, fmt.Sprintf("creates placeholders in app/dist (useful for go:embed)"))
var FlagClean = flag.BoolP("clean", "", false, fmt.Sprintf("cleans project"))
var FlagDev = flag.BoolP("dev", "d", false, fmt.Sprintf("starts dev mode"))
var FlagBuild = flag.BoolP("build", "b", false, fmt.Sprintf("builds project"))
var FlagHooks = flag.BoolP("hooks", "", false, fmt.Sprintf("adds git hooks"))
var FlagConfigure = flag.BoolP("configure", "", false, fmt.Sprintf("configures project by installing necessary binaries under \"./.gen\""))
var FlagPlatform = flag.StringP("platform", "", "", fmt.Sprintf("sets the platform, accepts either \"Linux/amd64\", \"Darwin/arm64\" or \"Darwin/amd64\""))
var FlagYes = flag.BoolP("yes", "y", false, fmt.Sprintf("confirms all binary promps silently"))
var FlagGo = flag.StringP("go", "", "go", fmt.Sprintf("sets the go binary, defaults to \"go\""))
var FlagAir = flag.StringP("air", "", filepath.Join(".gen", "air", "air"), fmt.Sprintf("sets the air binary, defaults to \".gen/air/air\""))
var FlagBun = flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"), fmt.Sprintf("sets the bun binary, defaults to \".gen/bun/bun\""))
var FlagSqlite = flag.StringP("sqlite", "", filepath.Join(".gen", "sqlite", "sqlite3"), fmt.Sprintf("sets the sqlite binary, defaults to \".gen/sqlite/sqlite3\""))

func (cli *Cli) OnStart() {
	if !cli.Parsed {
		flag.Parse()
		cli.Parsed = true
	}

	if *FlagHelp {
		cli.OnHelp()
		os.Exit(0)
	}

	if *FlagVersion {
		cli.OnVersion()
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		cli.OnCreateProject(*FlagCreateProject)
		os.Exit(0)
	}

	if *FlagAdd != "" {
		cli.OnAddFeature(*FlagAdd)
		os.Exit(0)
	}

	if *FlagTest {
		cli.OnTest()
		os.Exit(0)
	}

	if *FlagPackage {
		cli.OnPackage()
		os.Exit(0)
	}

	if *FlagPackageWatch {
		cli.OnPackageWatch()
		os.Exit(0)
	}

	if *FlagCheck {
		cli.OnCheck()
		os.Exit(0)
	}

	if *FlagUpdate {
		cli.OnUpdate()
		os.Exit(0)
	}

	if *FlagInstall {
		cli.OnInstall()
		os.Exit(0)
	}

	if *FlagFormat {
		cli.OnFormat()
		os.Exit(0)
	}

	if *FlagTouch {
		cli.OnTouch()
		os.Exit(0)
	}

	if *FlagClean {
		cli.OnClean()
		os.Exit(0)
	}

	if *FlagDev {
		cli.OnDev()
		os.Exit(0)
	}

	if *FlagBuild {
		cli.OnBuild()
		os.Exit(0)
	}

	if *FlagHooks {
		cli.OnHooks()
		os.Exit(0)
	}

	if *FlagConfigure {
		cli.OnConfigure()
		os.Exit(0)
	}

	cli.OnMenu()
}

func (cli *Cli) OnMenu() {
	err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("Frizzante", pterm.FgCyan.ToStyle())).Render()
	if err != nil {
		cli.Fatal(err)
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
	}

	result, showError := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Pick an option")
	if showError != nil {
		cli.Fatal(showError)
	}

	if result == "Help" {
		*FlagHelp = true
		cli.OnStart()
		return
	}

	if result == "Version" {
		*FlagVersion = true
		cli.OnStart()
		return
	}

	if result == "Create Project" {
		projectName, projectNameError := pterm.DefaultInteractiveTextInput.Show("Give the project name")
		if projectNameError != nil {
			cli.Fatal(projectNameError)
		}
		*FlagCreateProject = projectName
		cli.OnStart()
		return
	}

	if result == "Add" {
		*FlagAdd = ":pick"
		cli.OnStart()
		return
	}

	if result == "Add?" {
		*FlagAdd = "?"
		cli.OnStart()
		return
	}

	if result == "Test" {
		*FlagTest = true
		cli.OnStart()
		return
	}

	if result == "Package" {
		*FlagPackage = true
		cli.OnStart()
		return
	}

	if result == "Package Watch" {
		*FlagPackageWatch = true
		cli.OnStart()
		return
	}

	if result == "Check" {
		*FlagCheck = true
		cli.OnStart()
		return
	}

	if result == "Update" {
		*FlagUpdate = true
		cli.OnStart()
		return
	}

	if result == "Install" {
		*FlagInstall = true
		cli.OnStart()
		return
	}

	if result == "Format" {
		*FlagFormat = true
		cli.OnStart()
		return
	}

	if result == "Touch" {
		*FlagTouch = true
		cli.OnStart()
		return
	}

	if result == "Clean" {
		*FlagClean = true
		cli.OnStart()
		return
	}

	if result == "Dev" {
		*FlagDev = true
		cli.OnStart()
		return
	}

	if result == "Build" {
		*FlagBuild = true
		cli.OnStart()
		return
	}

	if result == "Hooks" {
		*FlagHooks = true
		cli.OnStart()
		return
	}

	if result == "Configure" {
		*FlagConfigure = true
		cli.OnStart()
		return
	}
}

func (cli *Cli) OnHelp() {
	flag.Usage()
}

func (cli *Cli) OnVersion() {
	var version string

	versionData, versionError := cli.Efs.ReadFile("version")
	if versionError != nil {
		cli.Fatal(versionError)
	}

	version = string(versionData)

	println(version)
}

func (cli *Cli) OnCreateProject(project string) {
	downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", project+".zip")
	if downloadError != nil {
		cli.Fatal(downloadError)
	}

	unzipError := files.UnzipFile(project+".zip", project+".tmp")
	if unzipError != nil {
		cli.Fatal(unzipError)
	}

	removeError := os.Remove(project + ".zip")
	if removeError != nil {
		cli.Fatal(removeError)
	}

	renameError := os.Rename(filepath.Join(project+".tmp", "frizzante-starter-main"), project)
	if renameError != nil {
		cli.Fatal(renameError)
	}

	removeAllError := os.RemoveAll(filepath.Join(project + ".tmp"))
	if removeAllError != nil {
		cli.Fatal(removeAllError)
	}

	os.Exit(0)
}

func (cli *Cli) OnAddFeature(features string) {
	if features == "?" {
		cli.ShowFeaturesInfo()
		return
	}

	if features == ":pick" {
		selectedFeatures, showError := pterm.
			DefaultInteractiveMultiselect.
			WithKeySelect(keys.Space).
			WithKeyConfirm(keys.Enter).
			WithOptions([]string{
				"Core",
				"Forms",
				"Links",
				"Air",
				"Bun",
				"Sqlite",
			}).
			WithFilter(false).
			Show("Pick a feature to add")

		if showError != nil {
			cli.Fatal(showError)
		}

		for _, selectedFeature := range selectedFeatures {
			cli.AddFeatureByName(selectedFeature)
		}
		return
	}

	splitFeatures := strings.Split(features, ",")

	for _, feature := range splitFeatures {
		cli.AddFeatureByName(feature)
	}
	return
}

func (cli *Cli) CopyFeatureDirectories(instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		name := instruction.FeatureName
		from := instruction.OriginDirectory
		to := instruction.DestinationDirectory

		if files.IsDirectory(to) {
			if !cli.Confirmf("It looks like feature `%s` already exists in this project, would you like to overwrite it?", name) {
				cli.Infof("skipping `%s`", name)
				return
			}

			removeAllError := os.RemoveAll(to)
			if removeAllError != nil {
				cli.Fatal(removeAllError)
			}
		}

		unixFrom := strings.ReplaceAll(from, "\\", "/")
		unixFromFileNames, readDirError := embeds.ReadDir(cli.Efs, unixFrom)
		if readDirError != nil {
			cli.Fatal(readDirError)
		}

		for _, unixFileName := range unixFromFileNames {
			fileName := to + strings.ReplaceAll(strings.TrimPrefix(unixFileName, unixFrom), "/", string(filepath.Separator))
			directoryName := filepath.Dir(fileName)

			if !files.IsDirectory(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if mkdirError != nil {
					cli.Fatal(mkdirError)
				}
			}

			file, openError := os.Create(fileName)
			if openError != nil {
				cli.Fatal(openError)
			}

			esfFile, esfOpenError := cli.Efs.Open(unixFileName)
			if esfOpenError != nil {
				cli.Fatal(esfOpenError)
			}

			_, copyError := io.Copy(file, esfFile)
			if copyError != nil {
				cli.Fatal(copyError)
			}
		}

		cli.Successf("adding `%s`", to)
	}
}

func (cli *Cli) AddFeatureByName(feature string) {
	if strings.ToLower(feature) == "core" {
		core := filepath.Join("app", "frizzante", "core")
		cli.CopyFeatureDirectories([]FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
		return
	}

	if strings.ToLower(feature) == "forms" {
		core := filepath.Join("app", "frizzante", "core")
		forms := filepath.Join("app", "frizzante", "forms")

		if !files.IsDirectory(core) {
			if cli.Confirmf("It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				cli.CopyFeatureDirectories([]FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		cli.CopyFeatureDirectories([]FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: forms, DestinationDirectory: forms}})
		return
	}

	if strings.ToLower(feature) == "links" {
		core := filepath.Join("app", "frizzante", "core")
		links := filepath.Join("app", "frizzante", "links")

		if !files.IsDirectory(core) {
			if cli.Confirmf("It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				cli.CopyFeatureDirectories([]FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		cli.CopyFeatureDirectories([]FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: links, DestinationDirectory: links}})
		return
	}

	if strings.ToLower(feature) == "bun" {
		directoryName := filepath.Join(".gen", "bun")

		platform := cli.Platform()

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
		} else {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
		}

		cli.Install("bun", url, directoryName)

		var fileName string

		if platform == PlatformDarwinArm64 {
			fileName = filepath.Join(directoryName, "bun-darwin-aarch64", "bun")
		} else if platform == PlatformDarwinAmd64 {
			fileName = filepath.Join(directoryName, "bun-darwin-x64", "bun")
		} else if platform == PlatformLinuxAmd64 {
			fileName = filepath.Join(directoryName, "bun-linux-x64", "bun")
		} else {
			fileName = filepath.Join(directoryName, "bun-linux-x64", "bun")
		}

		if files.IsFile(fileName) {
			renameError := os.Rename(fileName, cli.Bun("."))
			if renameError != nil {
				cli.Fatal(renameError)
			}
		}

		if files.IsDirectory(fileName) {
			removeError := os.RemoveAll(filepath.Dir(fileName))
			if removeError != nil {
				cli.Fatal(removeError)
			}
		}

		return
	}

	if strings.ToLower(feature) == "sqlite" {
		directoryName := filepath.Join(".gen", "sqlite")

		platform := cli.Platform()

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://www.sqlite.org/2025/sqlite-tools-osx-arm64-3500300.zip"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://www.sqlite.org/2025/sqlite-tools-osx-x64-3500300.zip"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://www.sqlite.org/2025/sqlite-tools-linux-x64-3500300.zip"
		} else {
			url = "https://www.sqlite.org/2025/sqlite-tools-linux-x64-3500300.zip"
		}

		cli.Install("sqlite", url, directoryName)

		return
	}

	if strings.ToLower(feature) == "air" {
		directoryName := filepath.Join(".gen", "air")

		platform := cli.Platform()

		var url string

		if platform == PlatformDarwinArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
		} else if platform == PlatformDarwinAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
		} else if platform == PlatformLinuxAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
		} else {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
		}

		cli.Install("air", url, directoryName)

		return
	}

	cli.Fatalf("unknown feature `%s`", feature)
}

func (cli *Cli) Cwd() string {
	wd, wdDir := os.Getwd()
	if wdDir != nil {
		cli.Fatal(wdDir)
	}
	return wd
}

func (cli *Cli) OnTest() {
	cli.OnPackage()
	test := exec.Command(cli.Go("."), "test")
	test.Env = append(os.Environ(), "CGO_ENABLED=1")
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err := test.Run()
	if err != nil {
		cli.Fatal(err)
	}
}

func (cli *Cli) OnHooks() {
	fileName := ".git/hooks/pre-commit"
	directoryName := filepath.Dir(fileName)
	if !files.IsDirectory(directoryName) {
		cli.Fatalf("directory `%s` not found", directoryName)
		return
	}

	if files.IsFile(fileName) {
		if !cli.Confirm("This git repository already defines a pre-commit script, would you like to overwrite it?") {
			pterm.Info.Println("pre-commit hook skipped")
			return
		}
		removeError := os.Remove(fileName)
		if removeError != nil {
			cli.Fatal(removeError)
		}
		cli.Success("pre-commit script overwritten")
	}

	err := os.WriteFile(fileName, []byte("make test"), os.ModePerm)
	if err != nil {
		cli.Fatal(err)
	}

	cli.Success("hooks added")
}

func (cli *Cli) OnTouch() {
	touch := func(fileName string) {
		directoryName := filepath.Dir(fileName)

		if !files.IsDirectory(directoryName) {
			mkdirAllError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirAllError != nil {
				cli.Fatal(mkdirAllError)
			}
		}

		file, openError := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if openError != nil {
			cli.Fatal(openError)
		}

		closeError := file.Close()
		if closeError != nil {
			cli.Fatal(closeError)
		}
	}

	mkdirError := os.MkdirAll("app/dist", os.ModePerm)
	if mkdirError != nil {
		cli.Fatal(mkdirError)
	}

	touch("app/dist/server.js")
	touch("app/dist/client/index.html")
}

func (cli *Cli) OnClean() {
	clean := exec.Command(cli.Go("."), "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	runError := clean.Run()
	if runError != nil {
		cli.Fatal(runError)
	}

	removeError := os.RemoveAll("app/dist")
	if removeError != nil {
		cli.Fatal(removeError)
	}

	removeError = os.RemoveAll("app/node_modules")
	if removeError != nil {
		cli.Fatal(removeError)
	}

	removeError = os.RemoveAll(".gen/tmp")
	if removeError != nil {
		cli.Fatal(removeError)
	}

	removeError = os.RemoveAll(".vite")
	if removeError != nil {
		cli.Fatal(removeError)
	}

	cli.OnTouch()

	cli.Success("project cleaned")
}

func (cli *Cli) OnFormat() {
	cli.OnTouch()

	gofmt := exec.Command(cli.Go("."), "fmt")
	gofmt.Env = append(os.Environ())
	gofmt.Stderr = os.Stderr
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	gofmtError := gofmt.Run()
	if gofmtError != nil {
		cli.Fatal(gofmtError)
	}

	prettier := exec.Command(cli.Bun("app"), "x", "prettier", "--write", ".")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		cli.Fatal(prettierError)
	}

	cli.Success("project formatted")
}

func (cli *Cli) OnUpdate() {
	cli.OnTouch()

	get := exec.Command(cli.Go("."), "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	getError := get.Run()
	if getError != nil {
		cli.Fatal(getError)
	}

	prettier := exec.Command(cli.Bun("app"), "update")
	prettier.Dir = "app"
	prettier.Env = append(os.Environ())
	prettier.Stderr = os.Stderr
	prettier.Stdout = os.Stdout
	prettier.Stdin = os.Stdin
	prettierError := prettier.Run()
	if prettierError != nil {
		cli.Fatal(prettierError)
	}

	cli.Success("project dependencies updated")
}

func (cli *Cli) OnInstall() {
	cli.OnTouch()

	tidy := exec.Command(cli.Go("."), "mod", "tidy")
	tidy.Env = append(os.Environ())
	tidy.Stderr = os.Stderr
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	tidyError := tidy.Run()
	if tidyError != nil {
		cli.Fatal(tidyError)
	}

	install := exec.Command(cli.Bun("app"), "install")
	install.Dir = "app"
	install.Env = append(os.Environ())
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	install.Stdin = os.Stdin
	installError := install.Run()
	if installError != nil {
		cli.Fatal(installError)
	}

	cli.Success("project dependencies installed")
}

func (cli *Cli) OnPackage() {
	cli.OnTouch()

	server := exec.Command(cli.Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Run()
	if serverError != nil {
		cli.Fatal(serverError)
	}

	client := exec.Command(cli.Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Run()
	if clientError != nil {
		cli.Fatal(clientError)
	}

	//node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite
	esbuild := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esbuild.Dir = "app"
	esbuild.Env = append(os.Environ())
	esbuild.Stderr = os.Stderr
	esbuild.Stdout = os.Stdout
	esbuild.Stdin = os.Stdin
	esbuildError := esbuild.Run()
	if esbuildError != nil {
		cli.Fatal(esbuildError)
	}

	cli.Success("project app package generated in app/dist")
}

func (cli *Cli) OnPackageWatch() {
	cli.OnTouch()

	server := exec.Command(cli.Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = "app"
	server.Env = append(os.Environ())
	server.Stderr = os.Stderr
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Start()
	if serverError != nil {
		cli.Fatalf("vite server watcher failed to launch\n%s", serverError)
	}
	cli.Success("vite server watcher launched")

	client := exec.Command(cli.Bun("app"), "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	client.Dir = "app"
	client.Env = append(os.Environ())
	client.Stderr = os.Stderr
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Start()
	if clientError != nil {
		cli.Fatalf("vite client watcher failed to launch\n%s", clientError)
	}
	cli.Success("vite client watcher launched")

	clientWaitError := client.Wait()
	if clientWaitError != nil {
		cli.Fatal(clientWaitError)
	}

	serverWaitError := server.Wait()
	if serverWaitError != nil {
		cli.Fatal(serverWaitError)
	}
}

func (cli *Cli) OnDev() {
	cli.OnTouch()

	mkdirError := os.MkdirAll(filepath.Join(".gen", "tmp"), os.ModePerm)
	if mkdirError != nil {
		cli.Fatal(mkdirError)
	}

	air := exec.Command(cli.Air("."))
	air.Env = append(os.Environ(), "DEV=1", "CGO_ENABLED=1")
	air.Stderr = os.Stderr
	air.Stdout = os.Stdout
	air.Stdin = os.Stdin
	airError := air.Start()
	if airError != nil {
		cli.Fatalf("air watcher faield to launch\n%s", airError)
	}
	cli.Success("air watcher launched")

	var group sync.WaitGroup

	group.Add(1)

	go func() { cli.OnPackageWatch() }()

	group.Wait()
	tidyWaitError := air.Wait()
	if tidyWaitError != nil {
		cli.Fatal(tidyWaitError)
	}
}

func (cli *Cli) OnBuild() {
	cli.OnPackage()

	build := exec.Command(cli.Go("."), "build", "-o=.gen/bin/app", ".")
	build.Env = append(os.Environ(), "CGO_ENABLED=1")
	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	buildError := build.Run()
	if buildError != nil {
		cli.Fatal(buildError)
	}
	cli.Success("project built into .gen/bin/app")
}

func (cli *Cli) OnCheck() {
	cli.OnTouch()

	eslint := exec.Command(cli.Bun("app"), "x", "eslint")
	eslint.Dir = "app"
	eslint.Env = append(os.Environ())
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	eslintError := eslint.Run()
	if eslintError != nil {
		cli.Fatal(eslintError)
	}

	svelteCheck := exec.Command(cli.Bun("app"), "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = "app"
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	svelteCheckError := svelteCheck.Run()
	if svelteCheckError != nil {
		cli.Fatal(svelteCheckError)
	}
}

func (cli *Cli) OnConfigure() {
	cli.OnAddFeature("bun,air")
	cli.OnInstall()
}

func (cli *Cli) Install(name string, url string, destination string) {
	if files.IsDirectory(destination) {
		if !cli.Confirmf("It looks like `%s` is already installed in `%s`, would you like to overwrite it?", name, destination) {
			cli.Infof("skipping `%s`", name)
			return
		}

		removeError := os.RemoveAll(destination)
		if removeError != nil {
			cli.Fatal(removeError)
		}
	}

	spinner, spinnerError := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(fmt.Sprintf("installing %s...", name))
	if spinnerError != nil {
		cli.Fatal(spinnerError)
	}
	defer func() {
		stopError := spinner.Stop()
		if stopError != nil {
			cli.Fatal(stopError)
		}
	}()

	if !strings.HasSuffix(url, ".zip") {
		downloadError := files.DownloadFile(url, filepath.Join(destination, name))
		if downloadError != nil {
			cli.Fatal(downloadError)
		}

		cli.Successf("%s installed in `%s`", name, destination)
		return
	}

	zipFileName := destination + ".zip"
	downloadError := files.DownloadFile(url, zipFileName)
	if downloadError != nil {
		cli.Fatal(downloadError)
	}
	defer func() {
		removeError := os.Remove(zipFileName)
		if removeError != nil {
			cli.Fatal(removeError)
		}
	}()

	unzipError := files.UnzipFile(zipFileName, destination)
	if unzipError != nil {
		cli.Fatal(unzipError)
	}

	cli.Successf("%s installed in `%s`", name, destination)
}

func (cli *Cli) ShowFeaturesInfo() {
	cli.Info(strings.Join([]string{
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
			"Sqlite",
			strings.Join([]string{
				"Adds sqlite to the project.",
				"",
				"Binaries will be dropped in `.gen/sqlite`.",
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
		cli.Fatal(tableError)
	}
}

type Platform uint

const PlatformLinuxAmd64 Platform = 0
const PlatformDarwinAmd64 Platform = 1
const PlatformDarwinArm64 Platform = 2

func (cli *Cli) Platform() Platform {
	var platform string

	if *FlagPlatform != "" {
		platform = *FlagPlatform
	} else {
		var platformError error
		platform, platformError = pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Linux/amd64",
				"Darwin/amd64",
				"Darwin/arm64",
			}).
			WithFilter(false).
			Show("Pick a platform")

		if platformError != nil {
			cli.Fatal(platformError)
		}
		*FlagPlatform = platform
	}

	if platform == "Darwin/arm64" {
		return PlatformDarwinArm64
	}

	if platform == "Darwin/amd64" {
		return PlatformDarwinAmd64
	}

	return PlatformLinuxAmd64
}

func (cli *Cli) Confirm(text string) bool {
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
		cli.Fatal(showError)
	}

	return yes
}

func (cli *Cli) Go(basepath string) string {
	var goBinary string

	if *FlagGo != "" {
		goBinary = *FlagGo
	} else {
		goBinary = cli.Go(".")
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
		cli.Fatal(pathError)
	}

	return path
}

func (cli *Cli) Air(basepath string) string {
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
		cli.Fatal(pathError)
	}

	return path
}

func (cli *Cli) Bun(basepath string) string {
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
		cli.Fatal(pathError)
	}

	return path
}

func (cli *Cli) Sqlite(basepath string) string {
	var sqlite string

	if *FlagSqlite != "" {
		sqlite = *FlagSqlite
	} else {
		sqlite = filepath.Join(".gen", "sqlite", "sqlite3")
	}

	if strings.HasPrefix(sqlite, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		sqlite = strings.Replace(sqlite, "~", dirname, 1)
		return sqlite
	}

	if !strings.Contains(sqlite, string(filepath.Separator)) {
		return sqlite
	}

	path, pathError := filepath.Rel(basepath, sqlite)
	if pathError != nil {
		cli.Fatal(pathError)
	}

	return path
}

func (cli *Cli) Confirmf(template string, vars ...any) bool {
	return cli.Confirm(fmt.Sprintf(template, vars...))
}

func (cli *Cli) Fatalf(template string, vars ...any) {
	pterm.Fatal.Printfln(template, vars...)
}

func (cli *Cli) Warningf(template string, vars ...any) {
	pterm.Warning.Printfln(template, vars...)
}

func (cli *Cli) Infof(template string, vars ...any) {
	pterm.Info.Printfln(template, vars...)
}

func (cli *Cli) Successf(template string, vars ...any) {
	pterm.Success.Printfln(template, vars...)
}

func (cli *Cli) Fatal(vars ...any) {
	pterm.Fatal.Println(vars...)
}

func (cli *Cli) Warning(vars ...any) {
	pterm.Warning.Println(vars...)
}

func (cli *Cli) Info(vars ...any) {
	pterm.Info.Println(vars...)
}

func (cli *Cli) Success(vars ...any) {
	pterm.Success.Println(vars...)
}

func (cli *Cli) Section(vars ...any) {
	pterm.DefaultSection.Println(vars...)
}
