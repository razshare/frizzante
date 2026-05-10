package actions

import (
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	if err = Package(PackageOptions{
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}
	extension := extensions.Find()
	spin := spinners.New("building binary")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"build", "-o=" + filepath.Join(".gen", "bin", "app"+extension), "."},
	}) {
		messages.Error("could not build go source code")
		return
	}
	spinners.Stop(spin)
	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"+extension))
	return
}
