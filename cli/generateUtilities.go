package cli

import (
	"embed"
	"github.com/pterm/pterm"
	"log"
	"os"
	"path/filepath"
)

func GenerateUtilities(efs embed.FS) {
	if !*FlagGenerateUtilities {
		return
	}

	var e error

	if "" == *FlagOut {
		*FlagOut, e = pterm.
			DefaultInteractiveTextInput.
			Show("Drop utilities in")
		if e != nil {
			log.Fatal(e)
		}
	}

	u := NewAotUtilities(efs)

	if e = u.CreateOnDisk(filepath.Join(*FlagOut)); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}
