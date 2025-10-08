package login

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/form"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

// View renders the login form
func View(client *_client.Client) {
	// Create initial empty form state
	initialState := form.State{
		Data:    LoginForm{},
		Errors:  make(map[string][]string),
		Valid:   true,
		Tainted: make(map[string]bool),
		Message: "",
		Pending: false,
	}

	send.View(client, views.View{
		Name: "Login",
		Props: Props{
			Form: initialState,
		},
	})
}
