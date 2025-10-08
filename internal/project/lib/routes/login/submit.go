package login

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/form"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

// Submit handles the login form submission
func Submit(client *_client.Client) {
	// Parse form data
	state := form.Parse(client, LoginForm{})

	// Validate the form
	state = form.Validate(state)

	// If validation failed, send errors back to client
	if !state.Valid {
		send.View(client, views.View{
			Name: "Login",
			Props: Props{
				Form: state,
			},
		})
		return
	}

	// Get the validated data
	loginData := state.Data.(LoginForm)

	// For now, just simulate success
	client.Config.InfoLog.Printf("Login attempt: %s", loginData.Email)

	// Redirect to home or send success response
	send.Navigate(client, "/")
}
