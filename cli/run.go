package cli

import (
	"flag"
	"github.com/pterm/pterm"
	"log"
)

func Run() {

	flag.Parse()

	if !*FlagProject &&
		!*FlagRender &&
		!*FlagUtilities {

		generator, showError := pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Render",
				"Project",
				"Utilities",
			}).
			Show("Generate")

		if showError != nil {
			log.Fatal(showError)
		}

		if "Render" == generator {
			*FlagRender = true
		}

		if "Project" == generator {
			*FlagProject = true
		}

		if "Utilities" == generator {
			*FlagUtilities = true
		}
	}

	Render()
	Project()
	Utilities()
}
