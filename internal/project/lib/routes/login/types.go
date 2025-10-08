package login

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/razshare/frizzante/internal/project/lib/core/form"
	"github.com/razshare/frizzante/internal/project/lib/core/types"
)

func init() {
	_ = types.Generate[Props]()
	_ = types.Generate[LoginForm]()
}

// Props are the properties passed to the Login view
type Props struct {
	Form form.State `json:"form"`
}

// LoginForm represents the login form data with validation rules
type LoginForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// Validate implements the form.Validatable interface
func (f LoginForm) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.Password, validation.Required, validation.Length(8, 100)),
	)
}
