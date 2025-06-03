package cli

import (
	"flag"
	"github.com/pterm/pterm"
	"log"
)

func Run() {

	flag.Parse()

	if !*FlagProject &&
		!*FlagRouter &&
		!*FlagUtilities {

		generator, showError := pterm.
			DefaultInteractiveSelect.
			WithOptions([]string{
				"Project",
				"Router",
				"Utilities",
			}).
			Show("Generate")

		if showError != nil {
			log.Fatal(showError)
		}

		if "Project" == generator {
			*FlagProject = true
		}

		if "Router" == generator {
			*FlagRouter = true
		}

		if "Utilities" == generator {
			*FlagUtilities = true
		}
	}

	Router()
	Project()
	Utilities()
}
