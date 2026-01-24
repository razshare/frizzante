package renders

import (
	_ "embed"
	"regexp"
)

var NoScript = regexp.MustCompile(`<script.*>.*</script>`)
