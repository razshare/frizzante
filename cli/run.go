package cli

import (
	"flag"
	"github.com/pterm/pterm"
	"log"
)

func Run() {

	flag.Parse()

	if !*FlagProject &&
		!*FlagUtilities {

		generator, showError := pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Project",
				"Utilities",
			}).
			Show("Generate")

		if showError != nil {
			log.Fatal(showError)
		}

		if "Project" == generator {
			*FlagProject = true
		}

		if "Utilities" == generator {
			*FlagUtilities = true
		}
	}

	Project()
	Utilities()
}
