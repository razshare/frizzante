package receive

import (
	"log"
	"os"
	"strconv"
)

var MaxFormSize int64 = 2097152

func init() {
	if value := os.Getenv("FRIZZANTE_MAX_FORM_SIZE"); value != "" {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(value, 10, 64); err != nil {
			log.Fatal(err)
			return
		}
		MaxFormSize = parsed
	}
}
