package receive

import (
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

var FormMetadataCache = map[reflect.Type]*FormMetadata{}

// Form reads the next multipart form or url encoded form message from the
// client and stores it in the value pointed to by value.
func Form(client *clients.Client, value any) bool {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return false
	}

	if client.Request.Form == nil && client.Request.MultipartForm == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return false
		}
	}

	reflection := reflect.ValueOf(value)

	if reflection.Kind() != reflect.Pointer {
		client.Config.ErrorLog.Println("form value must be a pointer", stack.Trace())
		return false
	}

	reflection = reflection.Elem()
	type_ := reflection.Type()

	for index := range reflection.NumField() {
		var ok bool
		var metadata *FormMetadata
		if metadata, ok = FormMetadataCache[type_]; !ok {
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

			if reflectionField.IsExported() {
				metadata = &FormMetadata{
					Key:       key,
					Exported:  true,
					Value:     reflectionValue,
					Reference: reflectionValue.Interface(),
				}
			} else {
				metadata = &FormMetadata{
					Key:      key,
					Exported: false,
					Value:    reflectionValue,
				}
			}
			FormMetadataCache[type_] = metadata
		}

		if !metadata.Exported {
			continue
		}

		var err error
		var pointer any
		switch metadata.Reference.(type) {
		case string:
			pointer = client.Request.Form.Get(metadata.Key)

		case []byte:
			pointer = []byte(client.Request.Form.Get(metadata.Key))

		case bool:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseBool(text); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
				return false
			}

		case []bool:
			entries := client.Request.Form[metadata.Key]
			local := make([]bool, len(entries))

			for jndex, entry := range entries {
				var parsed bool
				if parsed, err = strconv.ParseBool(entry); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
					return false
				}
				local[jndex] = parsed
			}
			pointer = local

		case uint:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, err = strconv.ParseUint(text, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
				return false
			}
			pointer = uint(tmp)

		case []uint:
			entries := client.Request.Form[metadata.Key]
			local := make([]uint, len(entries))

			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
					return false
				}
				local[jndex] = uint(tmp)
			}
			pointer = local

		case uint32:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp uint64
			if tmp, err = strconv.ParseUint(text, 10, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
				return false
			}
			pointer = uint32(tmp)

		case []uint32:
			entries := client.Request.Form[metadata.Key]
			local := make([]uint32, len(entries))

			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 32); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
					return false
				}
				local[jndex] = uint32(tmp)
			}
			pointer = local

		case uint64:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseUint(text, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
				return false
			}

		case []uint64:
			entries := client.Request.Form[metadata.Key]
			local := make([]uint64, len(entries))

			for jndex, entry := range entries {
				var tmp uint64
				if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local

		case int:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, err = strconv.ParseInt(text, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
				return false
			}
			pointer = int(tmp)

		case []int:
			entries := client.Request.Form[metadata.Key]
			local := make([]int, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
					return false
				}
				local[jndex] = int(tmp)
			}
			pointer = local

		case int32:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp int64
			if tmp, err = strconv.ParseInt(text, 10, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
				return false
			}
			pointer = int32(tmp)

		case []int32:
			entries := client.Request.Form[metadata.Key]
			local := make([]int32, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 32); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
					return false
				}
				local[jndex] = int32(tmp)
			}
			pointer = local

		case int64:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseInt(text, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
				return false
			}

		case []int64:
			entries := client.Request.Form[metadata.Key]
			local := make([]int64, len(entries))

			for jndex, entry := range entries {
				var tmp int64
				if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local

		case float32:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			var tmp float64
			if tmp, err = strconv.ParseFloat(text, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
				return false
			}
			pointer = float32(tmp)

		case []float32:
			entries := client.Request.Form[metadata.Key]
			local := make([]float32, len(entries))

			for jndex, entry := range entries {
				var tmp float64
				if tmp, err = strconv.ParseFloat(entry, 32); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
					return false
				}
				local[jndex] = float32(tmp)
			}
			pointer = local

		case float64:
			text := client.Request.Form.Get(metadata.Key)
			if text == "" {
				continue
			}
			if pointer, err = strconv.ParseFloat(text, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
				return false
			}

		case []float64:
			entries := client.Request.Form[metadata.Key]
			local := make([]float64, len(entries))

			for jndex, entry := range entries {
				var tmp float64
				if tmp, err = strconv.ParseFloat(entry, 64); err != nil {
					client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
					return false
				}
				local[jndex] = tmp
			}
			pointer = local

		case *multipart.FileHeader:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				pointer = headers[0]
			}
		case multipart.FileHeader:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				pointer = *headers[0]
			}
		case []multipart.FileHeader:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				locals := make([]multipart.FileHeader, len(headers))
				for jndex, header := range headers {
					locals[jndex] = *header
				}
				pointer = locals
			}
		case []*multipart.FileHeader:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				pointer = headers
			}
		case multipart.File:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				header := *headers[0]
				var file multipart.File
				if file, err = header.Open(); err != nil {
					client.Config.ErrorLog.Println(err, stack.Trace())
					return false
				}
				pointer = file
			}
		case []multipart.File:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				locals := make([]multipart.File, len(headers))
				for jndex, header := range headers {
					var file multipart.File
					if file, err = header.Open(); err != nil {
						client.Config.ErrorLog.Println(err, stack.Trace())
						return false
					}
					locals[jndex] = file
				}
				pointer = locals
			}
		case []*multipart.File:
			if headers := client.Request.MultipartForm.File[metadata.Key]; len(headers) > 0 {
				locals := make([]*multipart.File, len(headers))
				for jndex, header := range headers {
					var file multipart.File
					if file, err = header.Open(); err != nil {
						client.Config.ErrorLog.Println(err, stack.Trace())
						return false
					}
					locals[jndex] = &file
				}
				pointer = locals
			}
		default:
			client.Config.ErrorLog.Println("unknown form value type for key "+metadata.Key, stack.Trace())
			return false
		}
		metadata.Value.Set(reflect.ValueOf(pointer))
	}

	return true
}
