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
var Extension = ""

func init() {
	if string(filepath.Separator) == "\\" {
		Extension = ".exe"
	}
}

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
		os.Exit(0)
	}

	events := &FeatureAddEvents{
		ConfirmOverwrite: func(feature string) bool {
			return cli.Confirm(
				fmt.Sprintf(
					"It looks like `%s` already exists, would you like to overwrite it?",
					feature,
				),
			)
		},
		ConfirmAddMissingDependency: func(feature string, dependency string) bool {
			return cli.Confirm(
				fmt.Sprintf(
					"It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?",
					feature,
					dependency,
				),
			)
		},
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
			}).
			WithFilter(false).
			Show("Pick a feature to add")

		if showError != nil {
			cli.Fatal(showError)
		}

		for _, selectedFeature := range selectedFeatures {
			cli.AddFeatureByName(selectedFeature, events)
		}
		os.Exit(0)
	}

	splitFeatures := strings.Split(features, ",")

	for _, feature := range splitFeatures {
		cli.AddFeatureByName(feature, events)
	}
	os.Exit(0)
}

func (cli *Cli) CopyFeatureDirectories(events *FeatureAddEvents, instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		from := instruction.From
		to := instruction.To

		if files.IsDirectory(to) {
			yes := events.ConfirmOverwrite(to)
			if !yes {
				cli.Infof("Skipping `%s`.", to)
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

		cli.Successf("Adding `%s`.", to)
	}
}

func (cli *Cli) CopyFeatureFiles(events *FeatureAddEvents, feature string, instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		from := instruction.From
		to := instruction.To

		if files.IsFile(to) {
			yes := events.ConfirmOverwrite(feature)
			if !yes {
				cli.Infof("Skipping feature `%s`.", feature)
				return
			}

			removeError := os.Remove(to)
			if removeError != nil {
				cli.Fatal(removeError)
			}
		}

		file, openError := os.Create(to)
		if openError != nil {
			cli.Fatal(openError)
		}

		esfFile, esfOpenError := cli.Efs.Open(strings.ReplaceAll(from, "\\", "/"))
		if esfOpenError != nil {
			cli.Fatal(esfOpenError)
		}

		_, copyError := io.Copy(file, esfFile)
		if copyError != nil {
			cli.Fatal(copyError)
		}

		cli.Successf("Feature `%s` added.", feature)
	}
}

func (cli *Cli) AddFeatureByName(feature string, events *FeatureAddEvents) {
	if strings.ToLower(feature) == "core" {
		core := filepath.Join("app", "frizzante", "core")
		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
		return
	}

	if strings.ToLower(feature) == "forms" {
		core := filepath.Join("app", "frizzante", "core")
		forms := filepath.Join("app", "frizzante", "forms")

		if !files.IsDirectory(core) {
			if events.ConfirmAddMissingDependency(feature, "Core") {
				cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
			}
		}

		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: forms, To: forms}})
		return
	}

	if strings.ToLower(feature) == "links" {
		core := filepath.Join("app", "frizzante", "core")
		links := filepath.Join("app", "frizzante", "links")

		if !files.IsDirectory(core) {
			if events.ConfirmAddMissingDependency(feature, "Core") {
				cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
			}
		}

		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: links, To: links}})
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
	test := exec.Command("go"+Extension, "test")
	test.Dir = cli.Cwd()
	test.Env = append(os.Environ(), "CGO_ENABLED=1")
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

	err := os.WriteFile(fileName, []byte("frizzante --test"), os.ModePerm)
	if err != nil {
		cli.Fatal(err)
	}

	cli.Success("hooks added")
}

func (cli *Cli) OnTouch() {
	touch := func(fileName string) {
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

	touch("app/dist/.gitkeep")
	touch("app/dist/server.js")
	touch("app/dist/client/index.html")
}

func (cli *Cli) OnClean() {
	clean := exec.Command("go"+Extension, "clean")
	clean.Dir = cli.Cwd()
	clean.Env = append(os.Environ())
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

	cli.Success("project cleaned")
}

func (cli *Cli) OnFormat() {
	cli.OnTouch()

	gofmt := exec.Command("go"+Extension, "fmt")
	gofmt.Dir = cli.Cwd()
	gofmt.Env = append(os.Environ())
	gofmt.Stdout = os.Stdout
	gofmt.Stdin = os.Stdin
	gofmtError := gofmt.Run()
	if gofmtError != nil {
		cli.Fatal(gofmtError)
	}

	prettier := exec.Command("bunx"+Extension, "prettier", "--write", ".")
	prettier.Dir = filepath.Join(cli.Cwd(), "app")
	prettier.Env = append(os.Environ())
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

	get := exec.Command("go"+Extension, "get", "-u", "./...")
	get.Dir = cli.Cwd()
	get.Env = append(os.Environ())
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	getError := get.Run()
	if getError != nil {
		cli.Fatal(getError)
	}

	prettier := exec.Command("bun"+Extension, "update")
	prettier.Dir = filepath.Join(cli.Cwd(), "app")
	prettier.Env = append(os.Environ())
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

	tidy := exec.Command("go"+Extension, "mod", "tidy")
	tidy.Dir = cli.Cwd()
	tidy.Env = append(os.Environ())
	tidy.Stdout = os.Stdout
	tidy.Stdin = os.Stdin
	tidyError := tidy.Run()
	if tidyError != nil {
		cli.Fatal(tidyError)
	}

	install := exec.Command("bun"+Extension, "install")
	install.Dir = filepath.Join(cli.Cwd(), "app")
	install.Env = append(os.Environ())
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

	server := exec.Command("bunx"+Extension, "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = filepath.Join(cli.Cwd(), "app")
	server.Env = append(os.Environ())
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Run()
	if serverError != nil {
		cli.Fatal(serverError)
	}

	client := exec.Command("bunx"+Extension, "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	client.Dir = filepath.Join(cli.Cwd(), "app")
	client.Env = append(os.Environ())
	client.Stdout = os.Stdout
	client.Stdin = os.Stdin
	clientError := client.Run()
	if clientError != nil {
		cli.Fatal(clientError)
	}

	//node_modules/.bin/esbuild dist/server.js --bundle --outfile=dist/server.js --format=cjs --allow-overwrite
	esbuild := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esbuild.Dir = filepath.Join(cli.Cwd(), "app")
	esbuild.Env = append(os.Environ())
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

	server := exec.Command("bunx"+Extension, "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=frizzante/core/scripts/server.ts")
	server.Dir = filepath.Join(cli.Cwd(), "app")
	server.Env = append(os.Environ())
	server.Stdout = os.Stdout
	server.Stdin = os.Stdin
	serverError := server.Start()
	if serverError != nil {
		cli.Fatalf("vite server watcher failed to launch\n%s", serverError)
	}
	cli.Success("vite server watcher launched")

	client := exec.Command("bunx"+Extension, "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	client.Dir = filepath.Join(cli.Cwd(), "app")
	client.Env = append(os.Environ())
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

	mkdirError := os.MkdirAll(".gen/tmp", os.ModePerm)
	if mkdirError != nil {
		cli.Fatal(mkdirError)
	}

	air := exec.Command("air" + Extension)
	air.Dir = cli.Cwd()
	air.Env = append(os.Environ(), "DEV=1", "CGO_ENABLED=1")
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
	cli.OnTouch()

	build := exec.Command("go"+Extension, "build", "-o .gen/bin/app", ".")
	build.Dir = cli.Cwd()
	build.Env = append(os.Environ(), "CGO_ENABLED=1")
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

	eslint := exec.Command("bunx"+Extension, "eslint")
	eslint.Dir = filepath.Join(cli.Cwd(), "app")
	eslint.Env = append(os.Environ())
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	eslintError := eslint.Start()
	if eslintError != nil {
		cli.Fatal(eslintError)
	}

	svelteCheck := exec.Command("bunx"+Extension, "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = filepath.Join(cli.Cwd(), "app")
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	svelteCheckError := svelteCheck.Run()
	if svelteCheckError != nil {
		cli.Fatal(svelteCheckError)
	}
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
				"The core of frizzante.",
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
				"A <Form> component which behaves like a <form> element",
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
				"A <Link> component which behaves like an <a> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
				"",
				"Source code will be dropped in `app/frizzante/links`.",
				"",
				"Requires `Core`.",
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

func (cli *Cli) Confirm(text string) bool {
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
