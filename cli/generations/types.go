package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/cli/tags"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Types(options TypesOptions) (err error) {
	if !files.IsDirectory("types") {
		return
	}
	var fileNames []string
	if fileNames, err = files.ReadDirectory("types"); err != nil {
		return
	}
	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, ".go") {
			args := []string{"run"}
			if len(options.Tags) > 0 {
				var parsedTags []string
				if parsedTags, err = tags.Parse(options.Tags); err != nil {
					return
				}
				tagsLocal := make([]string, 0)
				for _, parsedTag := range parsedTags {
					tagsLocal = append(tagsLocal, parsedTag)
				}
				args = append(args, fmt.Sprintf("-tags=%s", strings.Join(tagsLocal, ",")))
			}
			args = append(args, fmt.Sprintf(".%stypes", string(filepath.Separator)))
			if !messages.Command(messages.CommandOptions{
				Environment: os.Environ(),
				Program:     options.Go,
				Args:        args,
			}) {
				err = errors.New("something went wrong while running types")
				return
			}
			messages.Success("types generated")
			messages.Tip(
				"## usage example\n",
				"<script lang=\"ts\">\n",
				"    import type { Props } from \"$gen/types/main/lib/routes/todos/props\"\n",
				"    let { items }: Props = $props()\n",
				"</script>",
			)
			break
		}
	}
	return
}
