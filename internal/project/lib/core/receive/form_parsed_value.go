package receive

import (
	"mime/multipart"
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormParsedValue a value associated with the given key, parses it according to the value's type,
// and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormParsedValue does not modify value.
//
// If value fails to parse, FormParsedValue returns false.
func FormParsedValue(client *clients.Client, key string, value any) bool {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return false
	}

	if client.Request.Form == nil && client.Request.MultipartForm == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
		}
	}

	var err error
	switch proxy := value.(type) {
	case *string:
		*proxy = client.Request.Form.Get(key)
		return true
	case *[]byte:
		*proxy = []byte(client.Request.Form.Get(key))
		return true
	case *bool:
		text := client.Request.Form.Get(key)
		if value, err = strconv.ParseBool(text); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
			return false
		}
		return true
	case *[]bool:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]bool, len(entries))

		for index, entry := range entries {
			var parsed bool
			if parsed, err = strconv.ParseBool(entry); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
				return false
			}
			local[index] = parsed
		}
		*proxy = local
	case *uint:
		text := client.Request.Form.Get(key)
		var tmp uint64
		if tmp, err = strconv.ParseUint(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
			return false
		}
		*proxy = uint(tmp)
		return true
	case *[]uint:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]uint, len(entries))

		for index, entry := range entries {
			var tmp uint64
			if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
				return false
			}
			local[index] = uint(tmp)
		}
		*proxy = local
	case *uint32:
		text := client.Request.Form.Get(key)
		var tmp int64
		if tmp, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
			return false
		}
		*proxy = uint32(tmp)
		return true
	case *[]uint32:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]uint32, len(entries))

		for index, entry := range entries {
			var tmp uint64
			if tmp, err = strconv.ParseUint(entry, 10, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
				return false
			}
			local[index] = uint32(tmp)
		}
		*proxy = local
	case *uint64:
		text := client.Request.Form.Get(key)
		if value, err = strconv.ParseUint(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
			return false
		}
		return true
	case *[]uint64:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]uint64, len(entries))

		for index, entry := range entries {
			var tmp uint64
			if tmp, err = strconv.ParseUint(entry, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
				return false
			}
			local[index] = tmp
		}
		*proxy = local
	case *int:
		text := client.Request.Form.Get(key)
		var tmp int64
		if tmp, err = strconv.ParseInt(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
			return false
		}
		*proxy = int(tmp)
		return true
	case *[]int:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]int, len(entries))

		for index, entry := range entries {
			var tmp int64
			if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
				return false
			}
			local[index] = int(tmp)
		}
		*proxy = local
	case *int32:
		text := client.Request.Form.Get(key)
		var tmp int64
		if tmp, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
			return false
		}
		*proxy = int32(tmp)
		return true
	case *[]int32:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]int32, len(entries))

		for index, entry := range entries {
			var tmp int64
			if tmp, err = strconv.ParseInt(entry, 10, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
				return false
			}
			local[index] = int32(tmp)
		}
		*proxy = local
	case *int64:
		text := client.Request.Form.Get(key)
		if value, err = strconv.ParseInt(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
			return false
		}
		return true
	case *[]int64:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]int64, len(entries))

		for index, entry := range entries {
			var tmp int64
			if tmp, err = strconv.ParseInt(entry, 10, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
				return false
			}
			local[index] = tmp
		}
		*proxy = local
	case *float32:
		text := client.Request.Form.Get(key)
		var tmp float64
		if tmp, err = strconv.ParseFloat(text, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
			return false
		}
		*proxy = float32(tmp)
		return true
	case *[]float32:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]float32, len(entries))

		for index, entry := range entries {
			var tmp float64
			if tmp, err = strconv.ParseFloat(entry, 32); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
				return false
			}
			local[index] = float32(tmp)
		}
		*proxy = local
	case *float64:
		text := client.Request.Form.Get(key)
		if value, err = strconv.ParseFloat(text, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
			return false
		}
		return true
	case *[]float64:
		var ok bool
		var entries []string
		if entries, ok = client.Request.Form[key]; !ok {
			return false
		}

		local := make([]float64, len(entries))

		for index, entry := range entries {
			var tmp float64
			if tmp, err = strconv.ParseFloat(entry, 64); err != nil {
				client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
				return false
			}
			local[index] = tmp
		}
		*proxy = local
	case *multipart.FileHeader:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			*proxy = *headers[0]
			return true
		}
	case *[]multipart.FileHeader:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			locals := make([]multipart.FileHeader, len(headers))
			for index, header := range headers {
				locals[index] = *header
			}
			*proxy = locals
			return true
		}
	case *[]*multipart.FileHeader:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			*proxy = headers
			return true
		}
	case *multipart.File:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			header := *headers[0]
			var file multipart.File
			if file, err = header.Open(); err != nil {
				client.Config.ErrorLog.Println(err, stack.Trace())
				return false
			}
			*proxy = file
			return true
		}
	case *[]multipart.File:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			locals := make([]multipart.File, len(headers))
			for index, header := range headers {
				var file multipart.File
				if file, err = header.Open(); err != nil {
					client.Config.ErrorLog.Println(err, stack.Trace())
					return false
				}
				locals[index] = file
			}
			*proxy = locals
			return true
		}
	case *[]*multipart.File:
		if headers := client.Request.MultipartForm.File[key]; len(headers) > 0 {
			locals := make([]*multipart.File, len(headers))
			for index, header := range headers {
				var file multipart.File
				if file, err = header.Open(); err != nil {
					client.Config.ErrorLog.Println(err, stack.Trace())
					return false
				}
				locals[index] = &file
			}
			*proxy = locals
			return true
		}
	case *any:
		*proxy = client.Request.Form.Get(key)
		return true
	}
	return false
}
