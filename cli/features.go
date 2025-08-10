package cli

import (
	"embed"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/table"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func OnAddFeature(efs embed.FS, features string) {
	if features == "?" {
		ShowFeaturesInfo()
		return
	}

	if features == ":pick" {
		selectedFeatures, showError := multiselect.Send("Pick a feature to add", []string{
			"Core",
			"Forms",
			"Links",
			"Air",
			"Bun",
			"Sqlc",
		})

		if showError != nil {
			Fatal(showError)
		}

		for _, selectedFeature := range selectedFeatures {
			AddFeatureByName(efs, selectedFeature)
		}
		return
	}

	splitFeatures := strings.Split(features, ",")

	for _, feature := range splitFeatures {
		AddFeatureByName(efs, feature)
	}
	return
}

func CopyFeatureDirectories(efs embed.FS, instructions []FeatureCopyInstruction) {
	for _, instruction := range instructions {
		name := instruction.FeatureName
		from := instruction.OriginDirectory
		to := instruction.DestinationDirectory

		if files.IsDirectory(to) {
			if !Confirmf("It looks like feature `%s` already exists in this project, would you like to overwrite it?", name) {
				Infof("skipping `%s`", name)
				return
			}

			removeAllError := os.RemoveAll(to)
			if removeAllError != nil {
				Fatal(removeAllError)
			}
		}

		unixFrom := strings.ReplaceAll(from, "\\", "/")
		unixFromFileNames, readDirError := embeds.ReadDirectory(efs, unixFrom)
		if readDirError != nil {
			Fatal(readDirError)
		}

		for _, unixFileName := range unixFromFileNames {
			fileName := to + strings.ReplaceAll(strings.TrimPrefix(unixFileName, unixFrom), "/", string(filepath.Separator))
			directoryName := filepath.Dir(fileName)

			if !files.IsDirectory(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if mkdirError != nil {
					Fatal(mkdirError)
				}
			}

			file, openError := os.Create(fileName)
			if openError != nil {
				Fatal(openError)
			}

			esfFile, esfOpenError := efs.Open(unixFileName)
			if esfOpenError != nil {
				Fatal(esfOpenError)
			}

			_, copyError := io.Copy(file, esfFile)
			if copyError != nil {
				Fatal(copyError)
			}
		}

		Successf("adding `%s`", to)
	}
}

func AddFeatureByName(efs embed.FS, feature string) {
	if strings.ToLower(feature) == "core" {
		core := filepath.Join("app", "frizzante", "core")
		CopyFeatureDirectories(efs, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
		return
	}

	if strings.ToLower(feature) == "forms" {
		core := filepath.Join("app", "frizzante", "core")
		forms := filepath.Join("app", "frizzante", "forms")

		if !files.IsDirectory(core) {
			if Confirmf("It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				CopyFeatureDirectories(efs, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		CopyFeatureDirectories(efs, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: forms, DestinationDirectory: forms}})
		return
	}

	if strings.ToLower(feature) == "links" {
		core := filepath.Join("app", "frizzante", "core")
		links := filepath.Join("app", "frizzante", "links")

		if !files.IsDirectory(core) {
			if Confirmf("It looks like you're missing the `%s` feature, which is required by the `%s` feature, would you like to add it?", feature, "Core") {
				CopyFeatureDirectories(efs, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: core, DestinationDirectory: core}})
			}
		}

		CopyFeatureDirectories(efs, []FeatureCopyInstruction{{FeatureName: feature, OriginDirectory: links, DestinationDirectory: links}})
		return
	}

	if strings.ToLower(feature) == "bun" {
		directoryName := filepath.Join(".gen", "bun")

		platform := Platform()

		var url string

		if platform == PlatformTypeDarwinArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-aarch64.zip"
		} else if platform == PlatformTypeDarwinAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-darwin-x64.zip"
		} else if platform == PlatformTypeLinuxArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-aarch64.zip"
		} else if platform == PlatformTypeLinuxAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-linux-x64.zip"
		} else if platform == PlatformTypeWindowsArm64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
		} else if platform == PlatformTypeWindowsAmd64 {
			url = "https://github.com/oven-sh/bun/releases/download/bun-v1.2.19/bun-windows-x64-baseline.zip"
		}

		Install("bun", url, directoryName)

		var fileName string

		if platform == PlatformTypeDarwinArm64 {
			fileName = filepath.Join(directoryName, "bun-darwin-aarch64", "bun")
		} else if platform == PlatformTypeDarwinAmd64 {
			fileName = filepath.Join(directoryName, "bun-darwin-x64", "bun")
		} else if platform == PlatformTypeLinuxArm64 {
			fileName = filepath.Join(directoryName, "bun-linux-aarch64", "bun")
		} else if platform == PlatformTypeLinuxAmd64 {
			fileName = filepath.Join(directoryName, "bun-linux-x64", "bun")
		} else if platform == PlatformTypeWindowsArm64 {
			fileName = filepath.Join(directoryName, "bun-windows-x64-baseline", "bun.exe")
		} else if platform == PlatformTypeWindowsAmd64 {
			fileName = filepath.Join(directoryName, "bun-windows-x64-baseline", "bun.exe")
		}

		if files.IsFile(fileName) {
			renameError := os.Rename(fileName, Bun("."))
			if renameError != nil {
				Fatal(renameError)
			}
		}

		installDirectory := filepath.Dir(fileName)
		if files.IsDirectory(directoryName) {
			removeError := os.RemoveAll(installDirectory)
			if removeError != nil {
				Fatal(removeError)
			}
		}

		return
	}

	if strings.ToLower(feature) == "air" {
		directoryName := filepath.Join(".gen", "air")

		platform := Platform()

		var url string

		if platform == PlatformTypeDarwinArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_arm64"
		} else if platform == PlatformTypeDarwinAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_darwin_amd64"
		} else if platform == PlatformTypeLinuxArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_arm64"
		} else if platform == PlatformTypeLinuxAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_linux_amd64"
		} else if platform == PlatformTypeWindowsArm64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_arm64.exe"
		} else if platform == PlatformTypeWindowsAmd64 {
			url = "https://github.com/air-verse/air/releases/download/v1.62.0/air_1.62.0_windows_amd64.exe"
		}

		Install("air", url, directoryName)

		return
	}

	if strings.ToLower(feature) == "sqlc" {
		directoryName := filepath.Join(".gen", "sqlc")

		platform := Platform()

		var url string

		if platform == PlatformTypeDarwinArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_arm64.zip"
		} else if platform == PlatformTypeDarwinAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_darwin_amd64.zip"
		} else if platform == PlatformTypeLinuxArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_arm64.zip"
		} else if platform == PlatformTypeLinuxAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_linux_amd64.zip"
		} else if platform == PlatformTypeWindowsArm64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		} else if platform == PlatformTypeWindowsAmd64 {
			url = "https://github.com/sqlc-dev/sqlc/releases/download/v1.29.0/sqlc_1.29.0_windows_amd64.zip"
		}

		Install("sqlc", url, directoryName)

		writeSchemaSql := true
		writeQueriesSql := true
		writeSqlcYaml := true

		if files.IsFile("schema.sql") {
			writeSchemaSql = Confirm("File `schema.sql` already exists, would you like to overwrite it?")
		}

		if writeSchemaSql {
			writeError := os.WriteFile("schema.sql", make([]byte, 0), os.ModePerm)
			if writeError != nil {
				Fatal(writeError)
			}
			Success("schema.sql created")
		}

		if files.IsFile("queries.sql") {
			writeQueriesSql = Confirm("File `queries.sql` already exists, would you like to overwrite it?")
		}

		if writeQueriesSql {
			writeError := os.WriteFile("queries.sql", make([]byte, 0), os.ModePerm)
			if writeError != nil {
				Fatal(writeError)
			}
			Success("queries.sql created")
		}

		if files.IsFile("sqlc.yaml") {
			writeSqlcYaml = Confirm("File `sqlc.yaml` already exists, would you like to overwrite it?")
		}

		if writeSqlcYaml {
			data, readError := efs.ReadFile("sqlc.yaml")
			if readError != nil {
				Fatal(readError)
			}

			readError = os.WriteFile("sqlc.yaml", data, os.ModePerm)
			if readError != nil {
				Fatal(readError)
			}
			Success("sqlc.yaml created")
		}

		return
	}

	Fatalf("unknown feature `%s`", feature)
}

func ShowFeaturesInfo() {
	messages.Info(strings.Join([]string{
		"You can use -a or --add in order ",
		"to add new features to the project.",
		"",
		"The value passed in must follow ",
		"the syntax: `-a{feature},{feature}`",
		"where {feature} is the name of the ",
		"feature.",
		"",
		"For example, `-acore,forms` will ",
		"generate the core and forms",
		"features of frizzante respectively ",
		"in `app/frizzante/core` and ",
		"`app/frizzante/forms`.",
		"",
		"Feature are not case-sensitive.",
		"",
		"You can also use -a:pick or ",
		"--add :pick to pick feature ",
		"interactively.",
	}, "\n"))

	println()

	table.Send(
		[]string{"Feature Name", "Description"},
		[][]string{
			{
				"Core",
				strings.Join([]string{
					"• Adds the core of frizzante.",
					"• A bundle of scripts and components that manage",
					"view rendering, view transitions, automatic state management,",
					"provides commonly used functions.",
					"",
					"• Source code will be dropped in `app/frizzante/core`.",
				}, "\n"),
			},
			{
				"Forms",
				strings.Join([]string{
					"• Adds a <Form> component which behaves like a <form> element",
					"with some additional features that facilitate",
					"the usage of web standards.",
					"",
					"• Source code will be dropped in `app/frizzante/forms`.",
					"",
					"• Requires `Core`.",
				}, "\n"),
			},
			{
				"Links",
				strings.Join([]string{
					"• Adds a <Link> component which behaves like an <a> element",
					"with some additional features that facilitate",
					"the usage of web standards.",
					"",
					"• Source code will be dropped in `app/frizzante/links`.",
					"",
					"• Requires `Core`.",
				}, "\n"),
			},
			{
				"Bun",
				strings.Join([]string{
					"• Adds bun to the project.",
					"",
					"• Binaries will be dropped in `.gen/bun`.",
					"",
					"• Bun is required for development mode.",
				}, "\n"),
			},
			{
				"Sqlc",
				strings.Join([]string{
					"• Adds sqlc to the project.",
					"",
					"• Binaries will be dropped in `.gen/sqlc`.",
				}, "\n"),
			},
		},
	)
}
