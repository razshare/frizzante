package receive

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var FormMetadataCache = map[reflect.Type][]*FormFieldMetadata{}
var FormMetadataCacheMutex sync.Mutex

// Form reads the next multipart form or url encoded form message from the
// client and stores it in the value pointed to by value.
func Form(request *http.Request, value any) (err error) {
	isMultipart := true
	if request.Form == nil && request.MultipartForm == nil {
		if err = request.ParseMultipartForm(FormMaxSzie); err != nil {
			if errors.Is(err, http.ErrNotMultipart) {
				isMultipart = false
				err = nil
			} else {
				return
			}
		}
	}
	reflection := reflect.ValueOf(value)
	if reflection.Kind() != reflect.Pointer {
		err = errors.New("form value must be a pointer")
		return
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
		var parseError error
		var pointer any
		switch metadata.Reference.(type) {
		case string:
			pointer = request.Form.Get(metadata.Key)
		case []byte:
			pointer = []byte(request.Form.Get(metadata.Key))
		case bool:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, parseError = strconv.ParseBool(text); parseError != nil {
				err = errors.New("form value is not a valid bool")
				return
			}
		case []bool:
			entries := request.Form[metadata.Key]
			local := make([]bool, len(entries))
			for jndex, entry := range entries {
				var parsed bool
				if parsed, parseError = strconv.ParseBool(entry); parseError != nil {
					err = errors.New("form value is not a valid bool")
					return
				}
				local[jndex] = parsed
			}
			pointer = local

		case uint:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, parseError = strconv.ParseUint(text, 10, 64); parseError != nil {
				err = errors.New("form value is not a valid uint")
				return
			}
			pointer = uint(tmp)
		case []uint:
			entries := request.Form[metadata.Key]
			local := make([]uint, len(entries))

			for jndex, entry := range entries {
				var tmp uint64
				if tmp, parseError = strconv.ParseUint(entry, 10, 64); parseError != nil {
					err = errors.New("form value is not a valid uint")
					return
				}
				local[jndex] = uint(tmp)
			}
			pointer = local
		case uint32:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, parseError = strconv.ParseUint(text, 10, 32); parseError != nil {
				err = errors.New("form value is not a valid uint32")
				return
			}
			pointer = uint32(tmp)
		case []uint32:
			entries := request.Form[metadata.Key]
			local := make([]uint32, len(entries))
			for jndex, entry := range entries {
				var tmp uint64
				if tmp, parseError = strconv.ParseUint(entry, 10, 32); parseError != nil {
					err = errors.New("form value is not a valid uint32")
					return
				}
				local[jndex] = uint32(tmp)
			}
			pointer = local
		case uint64:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, parseError = strconv.ParseUint(text, 10, 64); parseError != nil {
				err = errors.New("form value is not a valid uint64")
				return
			}
		case []uint64:
			entries := request.Form[metadata.Key]
			local := make([]uint64, len(entries))
			for jndex, entry := range entries {
				var tmp uint64
				if tmp, parseError = strconv.ParseUint(entry, 10, 64); parseError != nil {
					err = errors.New("form value is not a valid uint64")
					return
				}
				local[jndex] = tmp
			}
			pointer = local
		case int:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, parseError = strconv.ParseInt(text, 10, 64); parseError != nil {
				err = errors.New("form value is not a valid int")
				return
			}
			pointer = int(tmp)
		case []int:
			entries := request.Form[metadata.Key]
			local := make([]int, len(entries))
			for jndex, entry := range entries {
				var tmp int64
				if tmp, parseError = strconv.ParseInt(entry, 10, 64); parseError != nil {
					err = errors.New("form value is not a valid int")
					return
				}
				local[jndex] = int(tmp)
			}
			pointer = local
		case int32:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, parseError = strconv.ParseInt(text, 10, 32); parseError != nil {
				err = errors.New("form value is not a valid int32")
				return
			}
			pointer = int32(tmp)
		case []int32:
			entries := request.Form[metadata.Key]
			local := make([]int32, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, parseError = strconv.ParseInt(entry, 10, 32); parseError != nil {
					err = errors.New("form value is not a valid int32")
					return
				}
				local[jndex] = int32(tmp)
			}
			pointer = local
		case int64:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, parseError = strconv.ParseInt(text, 10, 64); parseError != nil {
				err = errors.New("form value is not a valid int64")
				return
			}
		case []int64:
			entries := request.Form[metadata.Key]
			local := make([]int64, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, parseError = strconv.ParseInt(entry, 10, 64); parseError != nil {
					err = errors.New("form value is not a valid int64")
					return
				}
				local[jndex] = tmp
			}
			pointer = local
		case float32:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp float64
			if tmp, parseError = strconv.ParseFloat(text, 32); parseError != nil {
				err = errors.New("form value is not a valid float32")
				return
			}
			pointer = float32(tmp)
		case []float32:
			entries := request.Form[metadata.Key]
			local := make([]float32, len(entries))
			for jndex, entry := range entries {
				var tmp float64
				if tmp, parseError = strconv.ParseFloat(entry, 32); parseError != nil {
					err = errors.New("form value is not a valid float32")
					return
				}
				local[jndex] = float32(tmp)
			}
			pointer = local
		case float64:
			text := request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, parseError = strconv.ParseFloat(text, 64); parseError != nil {
				err = errors.New("form value is not a valid float64")
				return
			}
		case []float64:
			entries := request.Form[metadata.Key]
			local := make([]float64, len(entries))
			for jndex, entry := range entries {
				var tmp float64
				if tmp, parseError = strconv.ParseFloat(entry, 64); parseError != nil {
					err = errors.New("form value is not a valid float64")
					return
				}
				local[jndex] = tmp
			}
			pointer = local
		case multipart.FileHeader:
			if !isMultipart {
				err = errors.New("could not parse file in form because it is not multipart")
				return
			}
			if headers := request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				pointer = *headers[0]
			}
		case []multipart.FileHeader:
			if !isMultipart {
				err = errors.New("could not parse file in form because it is not multipart")
				return
			}
			if headers := request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				locals := make([]multipart.FileHeader, len(headers))
				for jndex, header := range headers {
					locals[jndex] = *header
				}
				pointer = locals
			}
		default:
			err = fmt.Errorf("unknown form value type for key %s", metadata.Key)
			return
		}
		if pointer != nil {
			reflection.Field(index).Set(reflect.ValueOf(pointer))
		}
	}
	return
}
