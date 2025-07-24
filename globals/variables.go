package globals

import "regexp"

var SessionKey = "session.json"
var NoScript = regexp.MustCompile(`<script.*>.*</script>`)
