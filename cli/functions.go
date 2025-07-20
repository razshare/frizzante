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
var FlagAdd = flag.StringP("add", "a", "", fmt.Sprintf("adds features, see  \"-a?\" or \"--add ?\" for more details"))

func (cli *Cli) OnStart() {
	flag.Parse()

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

	events := &FeatureAddEvents{
		ConfirmOverwrite: func(feature string) bool {
			return Confirm(
				fmt.Sprintf(
					"It looks like `%s` already exists, would you like to overwrite it?",
					feature,
				),
			)
		},
		ConfirmAddMissingDependency: func(feature string, dependency string) bool {
			return Confirm(
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
				"Form",
				"Link",
			}).
			WithFilter(false).
			Show("Pick a feature to add")

		if showError != nil {
			log.Fatal(showError)
		}

		for _, selectedFeature := range selectedFeatures {
			cli.AddFeatureByName(selectedFeature, events)
		}
		os.Exit(0)
	}

	splitFeatures := strings.Split(features, ",")

	progress, progressError := pterm.
		DefaultProgressbar.
		WithTotal(len(splitFeatures)).
		WithTitle("Adding features").
		Start()

	if progressError != nil {
		log.Fatal(progressError)
	}

	for _, feature := range splitFeatures {
		cli.AddFeatureByName(feature, events)
		progress.Increment()
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
				pterm.Info.Printfln("Skipping `%s`.", to)
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

		pterm.Success.Printfln("Adding `%s`.", to)
	}
}

func (cli *Cli) CopyFeatureFiles(events *FeatureAddEvents, feature string, instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		from := instruction.From
		to := instruction.To

		if files.IsFile(to) {
			yes := events.ConfirmOverwrite(feature)
			if !yes {
				pterm.Info.Printfln("Skipping feature `%s`.", feature)
				return
			}

			removeError := os.Remove(to)
			if removeError != nil {
				log.Fatal(removeError)
			}
		}

		file, openError := os.Create(to)
		if openError != nil {
			log.Fatal(openError)
		}

		esfFile, esfOpenError := cli.Efs.Open(strings.ReplaceAll(from, "\\", "/"))
		if esfOpenError != nil {
			log.Fatal(esfOpenError)
		}

		_, copyError := io.Copy(file, esfFile)
		if copyError != nil {
			log.Fatal(copyError)
		}

		pterm.Success.Printfln("Feature `%s` added.", feature)
	}
}

func (cli *Cli) AddFeatureByName(feature string, events *FeatureAddEvents) {
	if strings.ToLower(feature) == "core" {
		core := filepath.Join("app", "frizzante", "core")
		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
		return
	}

	if strings.ToLower(feature) == "form" {
		core := filepath.Join("app", "frizzante", "core")
		form := filepath.Join("app", "frizzante", "form")

		if !files.IsDirectory(core) {
			if events.ConfirmAddMissingDependency(feature, "Core") {
				cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
			}
		}

		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: form, To: form}})
		return
	}

	if strings.ToLower(feature) == "link" {
		core := filepath.Join("app", "frizzante", "core")
		link := filepath.Join("app", "frizzante", "link")

		if !files.IsDirectory(core) {
			if events.ConfirmAddMissingDependency(feature, "Core") {
				cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: core, To: core}})
			}
		}

		cli.CopyFeatureDirectories(events, []FeatureCopyInstruction{{From: link, To: link}})
		return
	}

	log.Fatalf("unknown feature `%s`", feature)
}

func ShowFeaturesInfo() {
	pterm.Info.Println(strings.Join([]string{
		"You can use -a or --add",
		"in order to add new features to the project.",
		"",
		"The value passed in must follow",
		"the syntax: `-a{feature}`",
		"where {feature} is the name of the feature.",
		"",
		"For example, `-acore` will generate the core",
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
				"",
				"Source code will be dropped in `app/frizzante/core`.",
			}, "\n"),
		},
		{
			"Form",
			strings.Join([]string{
				"A <Form> component which behaves like a <form> element",
				"with some additional features that facilitate",
				"the usage of web standards.",
				"",
				"Source code will be dropped in `app/frizzante/form`.",
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
				"Source code will be dropped in `app/frizzante/link`.",
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
