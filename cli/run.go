package cli

import (
	"flag"
	"github.com/pterm/pterm"
	"log"
)

func Run() {
	flag.Parse()

	if !*FlagDevelop &&
		!*FlagCreateProject &&
		!*FlagGenerateUtilities {
		request, showError := pterm.
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

		if "Create Project" == request {
			*FlagCreateProject = true
		}

		if "Generate Utilities" == request {
			*FlagGenerateUtilities = true
		}

		if "Develop" == request {
			*FlagDevelop = true
		}
	}

	CreateProject()
	GenerateUtilities()
	Develop()
}
