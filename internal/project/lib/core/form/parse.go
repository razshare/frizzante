package form

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	_client "main/lib/core/client"
	"main/lib/core/receive"
	"main/lib/core/stack"
)

// Parse parses form data from the request and populates the provided struct.
//
// It returns a State containing the parsed data and any parsing errors.
//
// The struct fields should have `form:"fieldName"` tags to map form fields to struct fields.
//
// Example:
//
//	type LoginForm struct {
//	    Email    string `form:"email"`
//	    Password string `form:"password"`
//	}
//
//	state := form.Parse(client, LoginForm{})
func Parse[T any](client *_client.Client, schema T) State {
	state := State{
		Data:    schema,
		Errors:  make(map[string][]string),
		Valid:   true,
		Tainted: make(map[string]bool),
		Message: "",
		Pending: false,
	}

	// Parse the form data
	if client.Request.Form == nil {
		if err := client.Request.ParseMultipartForm(receive.MaxFormSize); err != nil {
			// Not a multipart form, try regular form parsing
			if err := client.Request.ParseForm(); err != nil {
				client.Config.ErrorLog.Println("failed to parse form", err, stack.Trace())
				state.Valid = false
				state.Message = "Failed to parse form data"
				return state
			}
		}
	}

	// Populate struct from form data
	populated, err := PopulateStruct(schema, client.Request.Form)
	if err != nil {
		client.Config.ErrorLog.Println("failed to populate struct", err, stack.Trace())
		state.Valid = false
		state.Message = "Failed to parse form data"
		return state
	}

	state.Data = populated

	// Track which fields were submitted (tainted)
	for key := range client.Request.Form {
		state.Tainted[key] = true
	}

	return state
}

// PopulateStruct populates a struct from url.Values using reflection.
//
// It looks for `form:"fieldName"` tags on struct fields.
func PopulateStruct[T any](schema T, values url.Values) (T, error) {
	val := reflect.ValueOf(&schema).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get form tag
		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			// Use field name if no tag specified
			formTag = strings.ToLower(fieldType.Name)
		}

		// Get value from form
		formValue := values.Get(formTag)
		if formValue == "" {
			continue
		}

		// Set field based on type
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(formValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(formValue, 10, 64); err == nil {
				field.SetInt(intVal)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(formValue, 10, 64); err == nil {
				field.SetUint(uintVal)
			}
		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(formValue, 64); err == nil {
				field.SetFloat(floatVal)
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(formValue); err == nil {
				field.SetBool(boolVal)
			}
		case reflect.Slice:
			// Handle multiple values for slices
			formValues := values[formTag]
			if len(formValues) > 0 {
				sliceVal := reflect.MakeSlice(field.Type(), len(formValues), len(formValues))
				for j, v := range formValues {
					elem := sliceVal.Index(j)
					if elem.Kind() == reflect.String {
						elem.SetString(v)
					}
				}
				field.Set(sliceVal)
			}
		}
	}

	return schema, nil
}

// ParseJSON parses form data from JSON request body.
//
// This is useful for API endpoints that accept JSON instead of form data.
func ParseJSON[T any](client *_client.Client, schema T) State {
	state := State{
		Data:    schema,
		Errors:  make(map[string][]string),
		Valid:   true,
		Tainted: make(map[string]bool),
		Message: "",
		Pending: false,
	}

	decoder := json.NewDecoder(client.Request.Body)
	if err := decoder.Decode(&schema); err != nil {
		client.Config.ErrorLog.Println("failed to parse JSON", err, stack.Trace())
		state.Valid = false
		state.Message = "Failed to parse JSON data"
		return state
	}

	state.Data = schema

	// Mark all fields as tainted for JSON requests
	val := reflect.ValueOf(schema)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			state.Tainted[jsonTag] = true
		} else {
			state.Tainted[strings.ToLower(fieldType.Name)] = true
		}
	}

	return state
}
