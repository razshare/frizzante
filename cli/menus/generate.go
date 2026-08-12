package menus

import (
	"fmt"
	"slices"

	"github.com/razshare/frizzante/v2/cli/actions"
	"github.com/razshare/frizzante/v2/cli/apps"
	"github.com/razshare/frizzante/v2/cli/generations"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/configs"
	"github.com/razshare/frizzante/v2/tui/search"
)

func init() {
	if files.IsFile("ask.md") {
		Generate.Items = append(Generate.Items, Item{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool {
				return slices.Contains([]string{"ask-docs"}, value)
			},
			Choice: search.Choice{Id: "ask docs", Description: "ask documentation"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ ask docs"))
				err = generations.AskDocs(generations.AskDocsOptions{})
				return
			},
		})
	}
	Generate.Items = append(Generate.Items, Item{
		Hidden: true,
		Choice: search.Choice{Id: "render generate menu"},
		Active: func(menu *Menu, app apps.App, value string, query []string) bool { return true },
		Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
			if *app.Strict {
				err = actions.Help(actions.HelpOptions{})
				return
			}
			for {
				if _, err = Render(menu, app, value, query, depth+1); err != nil {
					return
				}
				if depth > 1 {
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
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air"))
				return generations.Air(generations.AirOptions{Air: *app.Air})
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "air-config" },
			Choice: search.Choice{Id: "air-config", Description: "air configuration file air.toml"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))
				if err = generations.AirConfig(generations.AirConfigOptions{
					Efs: app.Efs,
				}); err != nil {
					return
				}
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "bun" },
			Choice: search.Choice{Id: "bun", Description: "bun binaries, a fast javascript all-in-one toolkit"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ bun"))
				err = generations.Bun(generations.BunOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "sqlc" },
			Choice: search.Choice{Id: "sqlc", Description: "sqlc binaries, a sql compiler written in go"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ sqlc"))
				err = generations.Sqlc(generations.SqlcOptions{Sqlc: *app.Sqlc})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "migration" },
			Choice: search.Choice{Id: "migration", Description: "migration file named using the current date and time"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
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
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "types" },
			Choice: search.Choice{Id: "types", Description: "typescript type definitions from go types"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ types"))
				err = generations.Types(generations.TypesOptions{
					Go:   *app.Go,
					Tags: *app.Tags,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "schema" },
			Choice: search.Choice{Id: "schema", Description: "go schema"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ schema"))
				err = generations.Schema(generations.SchemaOptions{
					Strict:   *app.Strict,
					Sqlc:     *app.Sqlc,
					SqlcYaml: *app.SqlcYaml,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "makefile" },
			Choice: search.Choice{Id: "makefile", Description: "makefile"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ makefile"))
				err = generations.Makefile(generations.MakefileOptions{
					Efs:    app.Efs,
					Strict: *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "snapshot" },
			Choice: search.Choice{Id: "snapshot", Description: "generates static pages, data and assets"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ snapshot"))
				err = generations.Snapshot(generations.SnapshotOptions{
					StaticsUrl: value,
					Strict:     *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "core" },
			Choice: search.Choice{Id: "core", Description: "core package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ core"))
				err = generations.Core(generations.CoreOptions{
					Efs:    app.Efs,
					Strict: *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "databases" },
			Choice: search.Choice{Id: "databases", Description: "databases package"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ databases"))
				err = generations.Databases(generations.DatabasesOptions{
					Efs:    app.Efs,
					Strict: *app.Strict,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App, value string, query []string) bool { return value == "icons" },
			Choice: search.Choice{Id: "icons", Description: "icons components"},
			Handle: func(menu *Menu, app apps.App, value string, query []string, depth int) (err error) {
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
	},
}
