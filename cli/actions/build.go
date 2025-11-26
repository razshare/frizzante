package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	if err = Package(PackageOptions{Bun: options.Bun, Prod: true}); err != nil {
		return
	}

	extension := extensions.Find()

	if len(options.Tags) > 0 {
		spin := spinners.Newf("building binary with tags %s", strings.Join(options.Tags, ","))
		go spinners.Start(spin)
		if messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     options.Go,
			Args:        []string{"build", "-tags=" + strings.Join(options.Tags, ","), "-o=" + filepath.Join(".gen", "bin", "app"+extension), "."},
		}) {
			spinners.Stop(spin)
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"+extension))
		} else {
			spinners.Stop(spin)
		}
	} else {
		spin := spinners.New("building binary")
		go spinners.Start(spin)
		if messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     options.Go,
			Args:        []string{"build", "-o=" + filepath.Join(".gen", "bin", "app"+extension), "."},
		}) {
			spinners.Stop(spin)
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"+extension))
		} else {
			spinners.Stop(spin)
		}
	}

	return
}
