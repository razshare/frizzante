package cli

import (
	"github.com/pterm/pterm"
	"github.com/razshare/frizzante/frz"
	"log"
	"os"
	"path/filepath"
)

func Utilities() {
	if !*FlagGenerate || !*FlagUtilities {
		return
	}

	var e error

	if "" == *FlagOut {
		*FlagOut, e = pterm.
			DefaultInteractiveTextInput.
			WithDelimiter("").
			Show("Drop utilities in")
		if e != nil {
			log.Fatal(e)
		}
	}

	u := frz.NewAotUtilities()

	if e = u.CreateOnDisk(filepath.Join(*FlagOut)); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}
