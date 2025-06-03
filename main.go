package main

import (
	"flag"
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/cli"
	"log"
)

func main() {
	flag.Parse()

	if !*cli.FlagProject &&
		!*cli.FlagRouter &&
		!*cli.FlagUtilities {

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
			*cli.FlagProject = true
		}

		if "Router" == generator {
			*cli.FlagRouter = true
		}

		if "Utilities" == generator {
			*cli.FlagUtilities = true
		}
	}

	cli.Router()
	cli.Project()
	cli.Utilities()
}
