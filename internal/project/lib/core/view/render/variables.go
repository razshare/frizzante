package render

import (
	_ "embed"
	"regexp"
)

//go:embed render.format
var RenderFormat string

//go:embed head.format
var HeadFormat string

//go:embed body.format
var BodyFormat string

//go:embed data.format
var DataFormat string

var NoScript = regexp.MustCompile(`<script.*>.*</script>`)
