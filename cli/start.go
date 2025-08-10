package cli

import (
	"embed"
	flag "github.com/spf13/pflag"
	"os"
)

func OnStart(efs embed.FS) {
	if !Parsed {
		flag.Parse()
		Parsed = true
	}

	if *FlagHelp {
		OnHelp()
		os.Exit(0)
	}

	if *FlagVersion {
		OnVersion(efs)
		os.Exit(0)
	}

	if *FlagCreateProject != "" {
		OnCreateProject(*FlagCreateProject)
		os.Exit(0)
	}

	if *FlagAdd != "" {
		OnAddFeature(efs, *FlagAdd)
		os.Exit(0)
	}

	if *FlagTest {
		OnTest()
		os.Exit(0)
	}

	if *FlagPackage {
		OnPackage()
		os.Exit(0)
	}

	if *FlagPackageWatch {
		OnPackageWatch()
		os.Exit(0)
	}

	if *FlagCheck {
		OnCheck()
		os.Exit(0)
	}

	if *FlagUpdate {
		OnUpdate()
		os.Exit(0)
	}

	if *FlagInstall {
		OnInstall()
		os.Exit(0)
	}

	if *FlagFormat {
		OnFormat()
		os.Exit(0)
	}

	if *FlagTouch {
		OnTouch()
		os.Exit(0)
	}

	if *FlagClean {
		OnClean()
		os.Exit(0)
	}

	if *FlagDev {
		OnDev()
		os.Exit(0)
	}

	if *FlagBuild {
		OnBuild()
		os.Exit(0)
	}

	if *FlagConfigure {
		OnConfigure(efs)
		os.Exit(0)
	}

	if *FlagSqlcGenerate {
		OnSqlcGenerate()
		os.Exit(0)
	}

	if *FlagCreateSqliteDatabase {
		OnSqliteDatabase(efs)
		os.Exit(0)
	}

	if *FlagWelcome {
		OnWelcome()
		os.Exit(0)
	}

	OnMenu(efs)
}
