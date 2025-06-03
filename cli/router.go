package cli

import (
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/frz"
	"log"
	"os"
	"path/filepath"
)

func Router() {
	if !*FlagGenerate || !*FlagRouter {
		return
	}

	var e error

	if "" == *FlagViews {
		*FlagViews, e = pterm.
			DefaultInteractiveTextInput.
			Show("Views are located in")
		if e != nil {
			log.Fatal(e)
		}
	}

	if "" == *FlagOut {
		*FlagOut, e = pterm.
			DefaultInteractiveTextInput.
			Show("Generate files in")
		if e != nil {
			log.Fatal(e)
		}
	}

	r := frz.NewAotRouter()

	if e = r.LoadViews(*FlagViews); e != nil {
		log.Fatal(e)
	}

	if e = r.CreateOnDisk(filepath.Join(*FlagOut)); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}
