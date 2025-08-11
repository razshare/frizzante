package globals

import "regexp"

var NoScript = regexp.MustCompile(`<script.*>.*</script>`)
