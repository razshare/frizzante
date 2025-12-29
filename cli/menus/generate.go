package menus

import (
	"fmt"

	"github.com/razshare/frizzante/cli/actions"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/search"
)

var Generate = Menu{
	Title: "generate",
	Items: []Item{
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "air" },
			Choice: search.Choice{Id: "air", Description: "generates air binaries, a live reload program for Go apps"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air"))
				return generate.Air(generate.AirOptions{Air: *app.Air})
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "air.toml" },
			Choice: search.Choice{Id: "air config", Description: "generates air configuration file air.toml"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ air config"))
				if err = generate.AirConfig(generate.AirConfigOptions{
					Efs:  app.Efs,
					Tags: *app.Tags,
				}); err != nil {
					return
				}
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "bun" },
			Choice: search.Choice{Id: "bun", Description: "generates bun binaries, a fast javascript all-in-one toolkit"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ bun"))
				err = generate.Bun(generate.BunOptions{Bun: *app.Bun})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "sqlc" },
			Choice: search.Choice{Id: "sqlc", Description: "generates sqlc binaries, a sql compiler written in go"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ sqlc"))
				err = generate.Sqlc(generate.SqlcOptions{Sqlc: *app.Sqlc})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "databases" },
			Choice: search.Choice{Id: "databases", Description: "generates databases package"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ databases"))
				err = generate.Databases(generate.DatabasesOptions{
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
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "migration" },
			Choice: search.Choice{Id: "migration", Description: "generates migration file named using the current date and time"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ migration"))
				err = generate.Migration(generate.MigrationOptions{
					Sqlc:     *app.Sqlc,
					Strict:   *app.Strict,
					SqlcYaml: *app.SqlcYaml,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "queries" },
			Choice: search.Choice{Id: "queries", Description: "generates go functions from your query file using sqlc"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ queries"))
				err = generate.Queries(generate.QueriesOptions{
					Strict:   *app.Strict,
					Sqlc:     *app.Sqlc,
					SqlcYaml: *app.SqlcYaml,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "core" },
			Choice: search.Choice{Id: "core", Description: "generates core package"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ core"))
				err = generate.Core(generate.CoreOptions{Efs: app.Efs})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "forms" },
			Choice: search.Choice{Id: "forms", Description: "generates forms components with error and pending handlers"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ forms"))
				err = generate.Forms(generate.FormsOptions{
					Strict: *app.Strict,
					Efs:    app.Efs})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "links" },
			Choice: search.Choice{Id: "links", Description: "generates links components with error and pending handlers"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ links"))
				err = generate.Links(generate.LinksOptions{
					Strict: *app.Strict,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "icons" },
			Choice: search.Choice{Id: "icons", Description: "generates icons components"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ icons"))
				err = generate.Icons(generate.IconsOptions{
					Strict: *app.Strict,
					Bun:    *app.Bun,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "security" },
			Choice: search.Choice{Id: "security", Description: "generates security package"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ security"))
				err = generate.Security(generate.SecurityOptions{
					Strict: *app.Strict,
					Efs:    app.Efs,
				})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "types" },
			Choice: search.Choice{Id: "types", Description: "generates typescript types from Go types"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ types"))
				err = generate.TypeDefinitions(generate.TypeDefinitionsOptions{Go: *app.Go})
				return
			},
		},
		{
			Active: func(menu *Menu, app apps.App) bool { return *app.Generate == "snapshot" },
			Choice: search.Choice{Id: "snapshot", Description: "generates static assets"},
			Handle: func(menu *Menu, app apps.App) (err error) {
				fmt.Print(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
				fmt.Println(configs.Styles.Menu.Render("generate ▷ snapshot"))
				return
			},
		},
		{
			Hidden: true,
			Choice: search.Choice{Id: "render generate menu"},
			Active: func(_ *Menu, app apps.App) bool { return true },
			Handle: func(menu *Menu, app apps.App) (err error) {
				if *app.Strict {
					err = actions.Help(actions.HelpOptions{})
					return
				}
				for {
					var id string
					if id, err = Render(menu, app); err != nil {
						return
					}

					if id == "" {
						return
					}
				}
			},
		},
	},
}
