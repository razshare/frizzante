package sandbox

import (
	"strings"
	
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/table"
)

func Preview() {
	messages.Info(strings.Join([]string{
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

	println()
	
	table.Send(
		[]string{"Feature Name", "Description"},
		[][]string{
			{
				"Core",
				strings.Join([]string{
					"Adds the core of frizzante.",
					"",
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
					"Adds a <Form> component which behaves like a <form> element",
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
					"Adds a <Link> component which behaves like an <a> element",
					"with some additional features that facilitate",
					"the usage of web standards.",
					"",
					"Source code will be dropped in `app/frizzante/links`.",
					"",
					"Requires `Core`.",
				}, "\n"),
			},
			{
				"Bun",
				strings.Join([]string{
					"Adds bun to the project.",
					"",
					"Binaries will be dropped in `.gen/bun`.",
					"",
					"Bun is required for development mode.",
				}, "\n"),
			},
			{
				"Sqlc",
				strings.Join([]string{
					"Adds sqlc to the project.",
					"",
					"Binaries will be dropped in `.gen/sqlc`.",
				}, "\n"),
			},
		},
	)
}
