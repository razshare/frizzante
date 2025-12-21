package menus

import "strings"

func NewActivationFunction(ids ...string) ActivationFunction {
	return func(query string) (value string, active bool) {
		for _, id := range ids {
			if strings.HasPrefix(query, id) {
				valueLocal := strings.TrimPrefix(query, id)

				// The reason we're doing this, is that we don't want the
				// order in which menu items are inserted into
				// the menu to matter for the activation.
				//
				// Take the following 2 use cases as an example.
				//
				// Available commands: "package-watch", "package"
				// Note the order.
				//
				// # Use case 1
				// User input: frizzante package-patch
				// Activated: "package-watch"    <====== correct.
				//
				// # Use case 2
				// User input: frizzante package
				// Activated: "package-watch"    <====== this is wrong,
				//				                         but if we don't this check
				//				                         it will behave like this.

				if valueLocal == "" || strings.HasPrefix(valueLocal, " ") {
					value = strings.TrimSpace(valueLocal)
					active = true
				}
				return
			}
		}
		return
	}
}
