package cli

import (
	"github.com/pterm/pterm"
	"log"
	"os"
	"path/filepath"
)

func Render() {
	if !*FlagGenerate || !*FlagRender {
		return
	}

	var e error

	if "" == *FlagViews {
		*FlagViews, e = pterm.
			DefaultInteractiveTextInput.
			Show("Pull views from")
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

	r := NewAotRouter()

	if e = r.LoadViews(*FlagViews); e != nil {
		log.Fatal(e)
	}

	if e = r.CreateOnDisk(filepath.Join(*FlagOut)); e != nil {
		log.Fatal(e)
	}

	os.Exit(0)
}
