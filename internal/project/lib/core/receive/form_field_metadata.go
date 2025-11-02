package receive

import "reflect"

type FormFieldMetadata struct {
	Key       string
	Value     reflect.Value
	Exported  bool
	Reference any
}
