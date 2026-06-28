package receive

import (
	"log"
	"os"
	"strconv"
)

var FormMaxSzie int64 = 2097152

func init() {
	if value := os.Getenv("FRIZZANTE_FORM_MAX_SIZE"); value != "" {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(value, 10, 64); err != nil {
			log.Fatal(err)
			return
		}
		FormMaxSzie = parsed
	}
}
