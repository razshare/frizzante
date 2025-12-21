package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/extensions"
	tags_ "github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	var tags []string
	if tags, err = tags_.Parse(options.Tags); err != nil {
		return
	}

	//if len(tags) == 0 && options.Interactive {
	//	if tags, err = tags_.Select([]search.Choice{
	//		{Id: "trace", Description: "enables tracing with stack.Trace()"},
	//		{Id: "no_js_runtime", Description: "disables the server-side JavaScript runtime"},
	//		{Id: "experimental_qjs_runtime", Description: "replaces goja with qjs"},
	//		{Id: "other", Description: "adds custom tags"},
	//	}); err != nil {
	//		return
	//	}
	//}

	if err = Package(PackageOptions{Bun: options.Bun, Production: true}); err != nil {
		return
	}

	extension := extensions.Find()

	if len(options.Tags) > 0 {
		spin := spinners.Newf("building binary with tags %s", strings.Join(tags, ","))
		go spinners.Start(spin)
		if messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     options.Go,
			Args:        []string{"build", "-tags=" + strings.Join(tags, ","), "-o=" + filepath.Join(".gen", "bin", "app"+extension), "."},
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
