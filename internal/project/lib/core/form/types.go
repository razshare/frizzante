package form

import (
)

func init() {
}

// State represents the complete form state that is sent between Go backend and Svelte frontend.
type State struct {
	Data    any                 `json:"data"`    // Form data (populated from request or defaults)
	Errors  map[string][]string `json:"errors"`  // Field-level validation errors
	Valid   bool                `json:"valid"`   // Whether the form passed validation
	Tainted map[string]bool     `json:"tainted"` // Which fields have been modified by the user
	Message string              `json:"message"` // Global form message (success/error)
	Pending bool                `json:"pending"` // Whether form is currently submitting (client-side only)
}

// Validatable is an interface that form structs can implement to define their validation rules.
//
// This uses ozzo-validation's Validate() method pattern.
type Validatable interface {
	Validate() error
}
