package menus

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
)

func init() {
	if err := LoadPlugins(&Generate, filepath.Join("plugins", "generate")); err != nil {
		messages.Fatal(err, stack.Trace())
		return
	}
	Generate.Items = append(Generate.Items, Item{
		Hidden: true,
		Choice: search.Choice{Id: "render generate menu"},
		Active: func(menu *Menu, app apps.App, value string, query []string) bool { return true },
		Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
			if *app.Strict {
				err = errors.New("unknown cli query")
				return
			}
			for {
				var id string
				if id, err = Render(menu, app, value, query); err != nil {
					return
				}
				if id == "" {
					return
				}
			}
		},
	})
}

var Generate = Menu{
	Title: "generate",
	Items: []Item{
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "air" },
			Choice: search.Choice{Id: "air", Description: "air binaries, a live reload program for Go apps"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air"))
				return generations.Air(generations.AirOptions{Air: *app.Air})
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "air-config" },
			Choice: search.Choice{Id: "air-config", Description: "air configuration file air.toml"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))
				if err = generations.AirConfig(generations.AirConfigOptions{
					Efs:  app.Efs,
					Tags: *app.Tags,
				}); err != nil {
					return
				}
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "bun" },
			Choice: search.Choice{Id: "bun", Description: "bun binaries, a fast javascript all-in-one toolkit"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ bun"))
				err = generations.Bun(generations.BunOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "sqlc" },
			Choice: search.Choice{Id: "sqlc", Description: "sqlc binaries, a sql compiler written in go"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ sqlc"))
				err = generations.Sqlc(generations.SqlcOptions{Sqlc: *app.Sqlc})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "databases" },
			Choice: search.Choice{Id: "databases", Description: "databases package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ databases"))
				err = generations.Databases(generations.DatabasesOptions{
					Go:     *app.Go,
					Sqlc:   *app.Sqlc,
					Type:   *app.DatabaseType,
					Strict: *app.Strict,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "migration" },
			Choice: search.Choice{Id: "migration", Description: "migration file named using the current date and time"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ migration"))
				err = generations.Migration(generations.MigrationOptions{
					Sqlc:     *app.Sqlc,
					Strict:   *app.Strict,
					SqlcYaml: *app.SqlcYaml,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "queries" },
			Choice: search.Choice{Id: "queries", Description: "go functions from your query file using sqlc"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ queries"))
				err = generations.Queries(generations.QueriesOptions{
					Strict:   *app.Strict,
					Sqlc:     *app.Sqlc,
					SqlcYaml: *app.SqlcYaml,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "snapshot" },
			Choice: search.Choice{Id: "snapshot", Description: "generates static pages, data and assets"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ snapshot"))
				staticsUrl := value
				var directoryName string
				if len(query) > 0 {
					directoryName = query[0]
				}
				err = generations.Snapshot(generations.SnapshotOptions{
					StaticsUrl:    staticsUrl,
					DirectoryName: directoryName,
					Strict:        *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "core" },
			Choice: search.Choice{Id: "core", Description: "core package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ core"))
				err = generations.Core(generations.CoreOptions{Efs: app.Efs})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "dev" },
			Choice: search.Choice{Id: "dev", Description: "dev package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ dev"))
				err = generations.Dev(generations.DevOptions{Efs: app.Efs})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "forms" },
			Choice: search.Choice{Id: "forms", Description: "forms components with error and pending handlers"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ forms"))
				err = generations.Forms(generations.FormsOptions{
					Strict: *app.Strict,
					Efs:    app.Efs})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "links" },
			Choice: search.Choice{Id: "links", Description: "links components with error and pending handlers"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ links"))
				err = generations.Links(generations.LinksOptions{
					Strict: *app.Strict,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "icons" },
			Choice: search.Choice{Id: "icons", Description: "icons components"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ icons"))
				err = generations.Icons(generations.IconsOptions{
					Strict: *app.Strict,
					Bun:    *app.Bun,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "security" },
			Choice: search.Choice{Id: "security", Description: "security package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ security"))
				err = generations.Security(generations.SecurityOptions{
					Strict: *app.Strict,
					Efs:    app.Efs,
				})
				return
			},
		},
	},
}
