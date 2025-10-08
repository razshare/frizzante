package form

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate validates a form struct using ozzo-validation.
//
// The struct must implement the Validatable interface (have a Validate() method).
//
// Returns a State with validation errors populated if validation fails.
//
// Example:
//
//	type LoginForm struct {
//	    Email    string `form:"email"`
//	    Password string `form:"password"`
//	}
//
//	func (f LoginForm) Validate() error {
//	    return validation.ValidateStruct(&f,
//	        validation.Field(&f.Email, validation.Required, is.Email),
//	        validation.Field(&f.Password, validation.Required, validation.Length(8, 100)),
//	    )
//	}
//
//	state := form.Parse(client, LoginForm{})
//	state = form.Validate(state)
func Validate(state State) State {
	// Check if data implements Validatable interface
	validatable, ok := state.Data.(Validatable)
	if !ok {
		// No validation defined, consider it valid
		return state
	}

	// Run validation
	err := validatable.Validate()
	if err == nil {
		// Validation passed
		state.Valid = true
		state.Errors = make(map[string][]string)
		return state
	}

	// Parse validation errors
	state.Valid = false
	state.Errors = ParseValidationErrors(err)

	return state
}

// ParseValidationErrors converts ozzo-validation errors into a map of field names to error messages.
func ParseValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)

	// Check if it's an ozzo-validation error
	if validationErrors, ok := err.(validation.Errors); ok {
		for field, fieldErr := range validationErrors {
			if fieldErr != nil {
				// Check if this field has nested errors
				if nestedErrors, ok := fieldErr.(validation.Errors); ok {
					// Handle nested struct validation errors
					for nestedField, nestedErr := range nestedErrors {
						key := field + "." + nestedField
						errors[key] = append(errors[key], nestedErr.Error())
					}
				} else {
					// Simple field error
					errors[field] = append(errors[field], fieldErr.Error())
				}
			}
		}
	} else {
		// Not an ozzo-validation error, add as general error
		errors["_form"] = append(errors["_form"], err.Error())
	}

	return errors
}

// ParseAndValidate is a convenience function that combines Parse and Validate.
//
// Example:
//
//	state := form.ParseAndValidate(client, LoginForm{})
//	if !state.Valid {
//	    send.View(client, view.View{
//	        Name: "Login",
//	        Props: map[string]any{"form": state},
//	    })
//	    return
//	}
func ParseAndValidate[T any](state State) State {
	return Validate(state)
}
