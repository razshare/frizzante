package cli

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	flag "github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var FlagHelp = flag.BoolP("help", "h", false, "shows this help document")
var FlagVersion = flag.BoolP("version", "v", false, "shows the Frizzante version used by this binary")
var FlagCreateProject = flag.StringP("create-project", "c", "", "creates a frizzante project")
var FlagAdd = flag.StringP("add", "a", "", fmt.Sprintf("adds features, see \"--add ?\" or \"-a?\" for more details"))

func (cli *Cli) Start() {
	flag.Parse()

	if *FlagVersion {
		cli.OnVersion()
		os.Exit(0)
	}

	if *FlagHelp {
		cli.OnHelp()
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

	cli.OnHelp()
}

func (cli *Cli) OnHelp() {
	flag.Usage()
}

func (cli *Cli) OnVersion() {
	var version string

	versionData, versionError := cli.Efs.ReadFile("version")
	if versionError != nil {
		log.Fatal(versionError)
	}

	version = string(versionData)

	println(version)
}

func (cli *Cli) OnCreateProject(project string) {
	downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", project+".zip")
	if downloadError != nil {
		log.Fatal(downloadError)
	}

	unzipError := files.UnzipFile(project+".zip", project+".tmp")
	if unzipError != nil {
		log.Fatal(unzipError)
	}

	removeError := os.Remove(project + ".zip")
	if removeError != nil {
		log.Fatal(removeError)
	}

	renameError := os.Rename(filepath.Join(project+".tmp", "frizzante-starter-main"), project)
	if renameError != nil {
		log.Fatal(renameError)
	}

	removeAllError := os.RemoveAll(filepath.Join(project + ".tmp"))
	if removeAllError != nil {
		log.Fatal(removeAllError)
	}

	os.Exit(0)
}

func (cli *Cli) OnAddFeature(features string) {
	if features == "?" {
		ShowFeaturesInfo()
		os.Exit(0)
	}

	for _, feature := range strings.Split(features, ",") {
		if feature == ":pick" {
			selectedFeatures, showError := pterm.
				DefaultInteractiveMultiselect.
				WithKeySelect(keys.Space).
				WithKeyConfirm(keys.Enter).
				WithOptions([]string{
					"Core",
					"Form",
					"Link",
				}).
				WithFilter(false).
				Show("Pick a feature to add")

			if showError != nil {
				log.Fatal(showError)
			}

			for _, selectedFeature := range selectedFeatures {
				cli.AddFeatureByName(selectedFeature)
			}
			continue
		}

		cli.AddFeatureByName(feature)
	}
	os.Exit(0)
}

func (cli *Cli) AddFeature(feature string, from string, to string) {
	if files.IsDirectory(to) {
		yes, _ := pterm.
			DefaultInteractiveConfirm.
			WithConfirmText("y").
			WithDefaultText("N").
			WithDefaultValue(false).
			Show(
				fmt.Sprintf(
					"It looks like you've already added feature `%s`, would you like to overwrite it?",
					feature,
				),
			)
		pterm.Println()

		if !yes {
			pterm.Info.Printfln("Skipping feature `%s`.", feature)
			return
		}

		removeAllError := os.RemoveAll(to)
		if removeAllError != nil {
			log.Fatal(removeAllError)
		}
	}

	unixFrom := strings.ReplaceAll(from, "\\", "/")
	unixFromFileNames, readDirError := embeds.ReadDir(cli.Efs, unixFrom)
	if readDirError != nil {
		log.Fatal(readDirError)
	}

	for _, unixFileName := range unixFromFileNames {
		fileName := to + strings.ReplaceAll(strings.TrimPrefix(unixFileName, unixFrom), "/", string(filepath.Separator))
		directoryName := filepath.Dir(fileName)

		if !files.IsDirectory(directoryName) {
			mkdirError := os.MkdirAll(directoryName, os.ModePerm)
			if mkdirError != nil {
				log.Fatal(mkdirError)
			}
		}

		file, openError := os.Create(fileName)
		if openError != nil {
			log.Fatal(openError)
		}

		esfFile, esfOpenError := cli.Efs.Open(unixFileName)
		if esfOpenError != nil {
			log.Fatal(esfOpenError)
		}

		_, copyError := io.Copy(file, esfFile)
		if copyError != nil {
			log.Fatal(copyError)
		}
	}

	pterm.Info.Printfln("Feature `%s` added.", feature)
}

func (cli *Cli) AddFeatureByName(feature string) {
	if feature == "core" {
		core := filepath.Join("app", "frizzante", "core")
		cli.AddFeature("Core", core, core)
		return
	}

	if strings.ToLower(feature) == "form" {
		core := filepath.Join("app", "frizzante", "core")
		form := filepath.Join("app", "frizzante", "form")

		if !files.IsDirectory(core) {
			if Confirm("It looks like you're missing the `Core` feature, which is required by the `Form` feature, would you like to add it?") {
				cli.AddFeature("Core", core, core)
			}
		}

		cli.AddFeature("Form", form, form)
		return
	}

	if strings.ToLower(feature) == "link" {
		core := filepath.Join("app", "frizzante", "core")
		link := filepath.Join("app", "frizzante", "link")

		if !files.IsDirectory(core) {
			if Confirm("It looks like you're missing the `Core` feature, which is required by the `Link` feature, would you like to add it?") {
				cli.AddFeature("Core", core, core)
			}
		}

		cli.AddFeature("Link", link, link)
		return
	}

	log.Fatalf("unknown feature `%s`", feature)

}

func ShowFeaturesInfo() {
	pterm.Info.Println(strings.Join([]string{
		"You can use --add or -a in combination with --feature or -f",
		"in order to add new features to the project.",
		"",
		"The value passed in must follow",
		"the syntax: `--add --feature {feature}`",
		"where {feature} is the name of the feature.",
		"",
		"Generated source code will be dropped in `app/frizzante`.",
		"",
		"For example, `--add --feature core` will generate the core",
		"features of frizzante in `app/frizzante/core`.",
		"",
		"NB: Feature names are not case-sensitive.",
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
			}, "\n"),
		},
		{
			"Form",
			strings.Join([]string{
				"A <Form> component which behaves like a <form> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
				"",
				"Requires `Core`.",
			}, "\n"),
		},
		{
			"Link",
			strings.Join([]string{
				"A <Link> component which behaves like an <a> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
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
		log.Fatal(tableError)
	}
}

func Confirm(text string) bool {
	yes, showError := pterm.
		DefaultInteractiveConfirm.
		WithConfirmText("Y").
		WithDefaultText("n").
		WithDefaultValue(true).
		Show(text)

	if showError != nil {
		log.Fatal(showError)
	}

	return yes
}
