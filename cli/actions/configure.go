package actions

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Configure(options ConfigureOptions) (err error) {
	if !files.IsFile(options.Bun) {
		if err = generations.Bun(generations.BunOptions{Bun: options.Bun}); err != nil {
			return
		}
	}
	if !files.IsFile(options.Air) {
		if err = generations.Air(generations.AirOptions{Air: options.Air}); err != nil {
			return
		}
	}
	if err = Install(InstallOptions{Go: options.Go, Bun: options.Bun}); err != nil {
		return
	}
	if err = generations.Types(generations.TypesOptions{Go: options.Go, Tags: options.Tags}); err != nil {
		return
	}
	if err = Package(PackageOptions{Go: options.Go, Bun: options.Bun}); err != nil {
		return
	}
	if err = PreBuild(PreBuildOptions{Go: options.Go, Tags: options.Tags}); err != nil {
		return
	}
	if err = BuildMigrate(BuildMigrateOptions{Go: options.Go, Tags: options.Tags, Output: options.Output}); err != nil {
		return
	}
	if !files.IsFile("source.sqlite") {
		messages.Info("database ./source.sqlite not found")
		spinner := spinners.New("migrating database")
		go spinners.Start(spinner)
		if !messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     filepath.Join(".gen", "bin", "migrate"),
		}) {
			spinners.Stop(spinner)
			err = errors.New("could not configure source.sqlite")
			return
		}
		spinners.Stop(spinner)
		messages.Success("database created and migrated into ./source.sqlite")
	}
	messages.Success("project configured")
	return
}
