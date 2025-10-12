package receive

import "reflect"

type FormMetadata struct {
	Key       string
	Value     reflect.Value
	Exported  bool
	Reference any
}
