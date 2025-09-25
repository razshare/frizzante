package guard

import "github.com/razshare/frizzante/internal/project/lib/core/client"

type HandlerFunc func(*client.Client)

// With wraps a handler with one or more guards
// Guards are executed in order, and the handler is only called if all guards pass
func With(handler HandlerFunc, guards ...Guard) HandlerFunc {
	return func(c *client.Client) {
		// Create a chain of guards
		// Start with the final handler
		finalAction := func() {
			handler(c)
		}

		// Build the chain in reverse order
		nextAction := finalAction
		for i := len(guards) - 1; i >= 0; i-- {
			// Capture the guard and next action in the closure properly
			g := guards[i]
			currentNext := nextAction

			// Create closure with captured values
			nextAction = func(guard Guard, next func()) func() {
				return func() {
					guard.Handler(c, next)
				}
			}(g, currentNext)
		}

		// Execute the chain
		nextAction()
	}
}
