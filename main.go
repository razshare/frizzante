package main

import (
	"atomicgo.dev/keyboard/keys"
	"embed"
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
var FlagVersion = flag.BoolP("version", "v", false, "shows the binary version and the project version")
var FlagCreateProject = flag.StringP("create-project", "c", "", fmt.Sprintf("creates a frizzante project"))
var FlagAdd = flag.StringP("add", "a", "", fmt.Sprintf("adds features, use \"--add ?\" or \"-a?\" for more details"))

//go:embed version
//go:embed app/frizzante
var efs embed.FS

func main() {
	flag.Parse()

	var version string

	versionData, versionError := efs.ReadFile("version")
	if versionError != nil {
		log.Fatal(versionError)
	}

	version = string(versionData)

	if *FlagVersion {
		println(version)
		os.Exit(0)
	}

	if *FlagHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		downloadError := files.DownloadFile("https://github.com/razshare/frizzante-starter/archive/refs/heads/main.zip", *FlagCreateProject+".zip")
		if downloadError != nil {
			log.Fatal(downloadError)
		}

		unzipError := files.UnzipFile(*FlagCreateProject+".zip", *FlagCreateProject+".tmp")
		if unzipError != nil {
			log.Fatal(unzipError)
		}

		removeError := os.Remove(*FlagCreateProject + ".zip")
		if removeError != nil {
			log.Fatal(removeError)
		}

		renameError := os.Rename(filepath.Join(*FlagCreateProject+".tmp", "frizzante-starter-main"), *FlagCreateProject)
		if renameError != nil {
			log.Fatal(renameError)
		}

		removeAllError := os.RemoveAll(filepath.Join(*FlagCreateProject + ".tmp"))
		if removeAllError != nil {
			log.Fatal(removeAllError)
		}

		os.Exit(0)
	}

	if *FlagAdd != "" {
		if *FlagAdd == "?" {
			SHowFeaturesInfo()
			os.Exit(0)
		}

		features := strings.Split(*FlagAdd, ",")

		for _, feature := range features {
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
					AddFeatureByName(strings.ToLower(strings.TrimSpace(selectedFeature)))
				}
				continue
			}

			AddFeatureByName(strings.ToLower(strings.TrimSpace(feature)))
		}
		os.Exit(0)
	}

	flag.Usage()
}

func AddFeature(feature string, from string, to string) {
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
	unixFromFileNames, readDirError := embeds.ReadDir(efs, unixFrom)
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

		esfFile, esfOpenError := efs.Open(unixFileName)
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

func MissingDependencyMessage(feature string, dependency string) string {
	return fmt.Sprintf(
		"It looks like you're missing the `%s` feature, which is required by `%s`, would you like to add it?",
		dependency,
		feature,
	)
}

func AddFeatureByName(feature string) {
	featureLower := strings.ToLower(feature)
	if featureLower == "core" {
		core := filepath.Join("app", "frizzante", "core")
		AddFeature("Core", core, core)
		return
	}

	if featureLower == "form" {
		core := filepath.Join("app", "frizzante", "core")
		form := filepath.Join("app", "frizzante", "form")

		if !files.IsDirectory(core) {
			if Confirm(MissingDependencyMessage(feature, "Core")) {
				AddFeature("Core", core, core)
			}
		}

		AddFeature("Form", form, form)
		return
	}

	if featureLower == "link" {
		core := filepath.Join("app", "frizzante", "core")
		link := filepath.Join("app", "frizzante", "link")

		if !files.IsDirectory(core) {
			if Confirm(MissingDependencyMessage(feature, "Core")) {
				AddFeature("Core", core, core)
			}
		}

		AddFeature("Link", link, link)
		return
	}

	log.Fatalf("unknown feature `%s`", feature)

}

func SHowFeaturesInfo() {
	pterm.Info.Println(strings.Join([]string{
		"You can use --add or -a",
		"in order to add new features",
		"to the project.",
		"",
		"The value passed in must follow",
		"the syntax: `--add {feature}`",
		"where {feature} is the name of the feature.",
		"",
		"Generated source code will be dropped in `app/frizzante`.",
		"",
		"For example, `--add core` will generate the core",
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
				"the usage of web standards.\n",
				"Requires `Core`.",
			}, "\n"),
		},
		{
			"Link",
			strings.Join([]string{
				"A <Link> component which behaves like an <a> element",
				"with some additional features that facilitate",
				"the usage of web standards.\n",
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
