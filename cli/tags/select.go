package tags

import (
	"slices"

	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/multiselect"
	"github.com/razshare/frizzante/tui/search"
)

func Select(choices []search.Choice) (tags []string, err error) {
	tags = make([]string, 0)
	var yes bool
	if yes, err = confirm.Send(false, "use build tags?"); err != nil {
		return
	}

	if yes {
		if tags, err = multiselect.Send(choices, "select build tags"); err != nil {
			return
		}

		if slices.Contains(tags, "other") {
			var answer string
			if answer, err = input.Send("add your custom tags separated by comma"); err != nil {
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
