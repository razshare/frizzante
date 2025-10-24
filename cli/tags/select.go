package tags

import (
	"slices"

	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_many"
)

func Select(choices []search.Choice) (tags []string, err error) {
	tags = make([]string, 0)
	var yes bool
	if yes, err = confirm.Send(false, "use build tags?"); err != nil {
		return
	}

	if yes {
		if tags, err = select_many.Send(choices, "select build tags"); err != nil {
			return
		}

		if slices.Contains(tags, "other") {
			var answer string
			if answer, err = inputs.Send("add your custom tags separated by comma"); err != nil {
				return
			}

			var parsedTags []string
			if parsedTags, err = Parse(answer); err != nil {
				return
			}

			tags = append(tags, parsedTags...)
		}
	}
	return
}
