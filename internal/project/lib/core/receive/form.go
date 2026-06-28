package receive

import (
	"errors"
	"mime/multipart"
	http_ "net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

var FormMetadataCache = map[reflect.Type][]*FormFieldMetadata{}
var FormMetadataCacheMutex sync.Mutex

// Form reads the next multipart form or url encoded form message from the
// client and stores it in the value pointed to by value.
func Form(http *scopes.Http, value any) bool {
	if http.WebSocket != nil {
		logs.Errorf(
			http,
			"receive.Form: web socket connections cannot parse forms\n%s",
			stack.Trace(),
		)
		return false
	}
	isMultipart := true
	if http.Request.Form == nil && http.Request.MultipartForm == nil {
		if err := http.Request.ParseMultipartForm(FormMaxSzie); err != nil {
			if errors.Is(err, http_.ErrNotMultipart) {
				isMultipart = false
			} else {
				logs.Errorf(
					http,
					"receive.Form: failed to parse multipart form: %v\n%s",
					err,
					stack.Trace(),
				)
				return false
			}
		}
	}
	reflection := reflect.ValueOf(value)
	if reflection.Kind() != reflect.Pointer {
		logs.Errorf(
			http,
			"receive.Form: form value must be a pointer\n%s",
			stack.Trace(),
		)
		return false
	}
	reflection = reflection.Elem()
	type_ := reflection.Type()
	var ok bool
	var cache []*FormFieldMetadata
	if cache, ok = FormMetadataCache[type_]; !ok {
		count := reflection.NumField()
		cache = make([]*FormFieldMetadata, count)
		for index := range count {
			reflectionField := type_.Field(index)
			var key string
			if tag := reflectionField.Tag.Get("form"); tag != "" {
				key = tag
			} else {
				if tag = reflectionField.Tag.Get("json"); tag != "" {
					key = tag
				} else {
					key = reflectionField.Name
				}
			}
			reflectionValue := reflection.Field(index)
			if reflectionValue.Kind() == reflect.Pointer {
				reflectionValue = reflectionValue.Elem()
			}
			var metadata *FormFieldMetadata
			if reflectionField.IsExported() {
				metadata = &FormFieldMetadata{
					Key:       key,
					Exported:  true,
					Value:     reflectionValue,
					Reference: reflectionValue.Interface(),
				}
			} else {
				metadata = &FormFieldMetadata{
					Key:      key,
					Exported: false,
					Value:    reflectionValue,
				}
			}
			cache[index] = metadata
		}
		FormMetadataCacheMutex.Lock()
		FormMetadataCache[type_] = cache
		FormMetadataCacheMutex.Unlock()
	}
	for index, metadata := range cache {
		if !metadata.Exported {
			continue
		}
		var err error
		var pointer any
		switch metadata.Reference.(type) {
		case string:
			pointer = http.Request.Form.Get(metadata.Key)
		case []byte:
			pointer = []byte(http.Request.Form.Get(metadata.Key))
		case bool:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseBool(text); err != nil {
				logs.Errorf(http, "form value is not a valid bool\n%s", stack.Trace())
				return false
			}
		case []bool:
			entries := http.Request.Form[metadata.Key]
			local := make([]bool, len(entries))
			for jndex, entry := range entries {
				var parsed bool
				if parsed, err = strconv.ParseBool(entry); err != nil {
					logs.Errorf(http, "form value is not a valid bool\n%s", stack.Trace())
					return false
				}
				local[jndex] = parsed
			}
			pointer = local

		case uint:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, err = strconv.ParseUint(text, 10, 64); err != nil {
				logs.Errorf(http, "form value is not a valid uint\n%s", stack.Trace())
				return false
			}
			pointer = uint(tmp)
		case []uint:
			entries := http.Request.Form[metadata.Key]
			local := make([]uint, len(entries))

			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
					logs.Errorf(http, "form value is not a valid uint\n%s", stack.Trace())
					return false
				}
				local[jndex] = uint(tmp)
			}
			pointer = local
		case uint32:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, err = strconv.ParseUint(text, 10, 32); err != nil {
				logs.Errorf(http, "form value is not a valid uint32\n%s", stack.Trace())
				return false
			}
			pointer = uint32(tmp)
		case []uint32:
			entries := http.Request.Form[metadata.Key]
			local := make([]uint32, len(entries))
			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 32); err != nil {
					logs.Errorf(http, "form value is not a valid uint32\n%s", stack.Trace())
					return false
				}
				local[jndex] = uint32(tmp)
			}
			pointer = local
		case uint64:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseUint(text, 10, 64); err != nil {
				logs.Errorf(http, "form value is not a valid uint64\n%s", stack.Trace())
				return false
			}
		case []uint64:
			entries := http.Request.Form[metadata.Key]
			local := make([]uint64, len(entries))
			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
					logs.Errorf(http, "form value is not a valid uint64\n%s", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local
		case int:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, err = strconv.ParseInt(text, 10, 64); err != nil {
				logs.Errorf(http, "form value is not a valid int\n%s", stack.Trace())
				return false
			}
			pointer = int(tmp)
		case []int:
			entries := http.Request.Form[metadata.Key]
			local := make([]int, len(entries))
			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
					logs.Errorf(http, "form value is not a valid int\n%s", stack.Trace())
					return false
				}
				local[jndex] = int(tmp)
			}
			pointer = local
		case int32:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, err = strconv.ParseInt(text, 10, 32); err != nil {
				logs.Errorf(http, "form value is not a valid int32\n%s", stack.Trace())
				return false
			}
			pointer = int32(tmp)
		case []int32:
			entries := http.Request.Form[metadata.Key]
			local := make([]int32, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 32); err != nil {
					logs.Errorf(http, "form value is not a valid int32\n%s", stack.Trace())
					return false
				}
				local[jndex] = int32(tmp)
			}
			pointer = local
		case int64:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseInt(text, 10, 64); err != nil {
				logs.Errorf(http, "form value is not a valid int64\n%s", stack.Trace())
				return false
			}
		case []int64:
			entries := http.Request.Form[metadata.Key]
			local := make([]int64, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
					logs.Errorf(http, "form value is not a valid int64\n%s", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local
		case float32:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp float64
			if tmp, err = strconv.ParseFloat(text, 32); err != nil {
				logs.Errorf(http, "form value is not a valid float32\n%s", stack.Trace())
				return false
			}
			pointer = float32(tmp)
		case []float32:
			entries := http.Request.Form[metadata.Key]
			local := make([]float32, len(entries))
			for jndex, entry := range entries {
				var tmp float64
				if tmp, err = strconv.ParseFloat(entry, 32); err != nil {
					logs.Errorf(http, "form value is not a valid float32\n%s", stack.Trace())
					return false
				}
				local[jndex] = float32(tmp)
			}
			pointer = local
		case float64:
			text := http.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseFloat(text, 64); err != nil {
				logs.Errorf(http, "form value is not a valid float64\n%s", stack.Trace())
				return false
			}
		case []float64:
			entries := http.Request.Form[metadata.Key]
			local := make([]float64, len(entries))
			for jndex, entry := range entries {
				var tmp float64
				if tmp, err = strconv.ParseFloat(entry, 64); err != nil {
					logs.Errorf(http, "form value is not a valid float64\n%s", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local
		case multipart.FileHeader:
			if !isMultipart {
				logs.Error(http, "could not parse file in form because it is not multipart")
				return false
			}
			if headers := http.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				pointer = *headers[0]
			}
		case []multipart.FileHeader:
			if !isMultipart {
				logs.Error(http, "could not parse file in form because it is not multipart")
				return false
			}
			if headers := http.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				locals := make([]multipart.FileHeader, len(headers))
				for jndex, header := range headers {
					locals[jndex] = *header
				}
				pointer = locals
			}
		default:
			logs.Errorf(http, "unknown form value type for key %s\n%s", metadata.Key, stack.Trace())
			return false
		}
		if pointer != nil {
			reflection.Field(index).Set(reflect.ValueOf(pointer))
		}
	}
	return true
}
