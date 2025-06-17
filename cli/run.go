package cli

import (
	"embed"
	"flag"
	"github.com/pterm/pterm"
	"log"
)

//go:embed .air.toml
//go:embed bin
var bin embed.FS

func Run() {
	flag.Parse()

	if !*FlagCreateProject &&
		!*FlagGenerateUtilities {

		generator, showError := pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Create Project",
				"Generate Utilities",
				"Develop",
			}).
			Show("Welcome to Frizzante, pick an option")

		if showError != nil {
			log.Fatal(showError)
		}

		if "Create Project" == generator {
			*FlagCreateProject = true
		}

		if "Generate Utilities" == generator {
			*FlagGenerateUtilities = true
		}

		if "Develop" == generator {
			*FlagDevelop = true
		}
	}

	Project()
	Utilities()
	*FlagDevelop = true
	Develop(bin)
}
