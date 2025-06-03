package main

import (
	"flag"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/cli/lib"
	"log"
)

func main() {
	flag.Parse()

	if !*lib.FlagProject &&
		!*lib.FlagRouter &&
		!*lib.FlagUtilities {

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
			*lib.FlagProject = true
		}

		if "Router" == generator {
			*lib.FlagRouter = true
		}

		if "Utilities" == generator {
			*lib.FlagUtilities = true
		}
	}

	lib.Router()
	lib.Project()
	lib.Utilities()
}
