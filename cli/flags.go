package cli

import "flag"

var FlagGenerate = flag.Bool("generate", true, "")
var FlagProject = flag.Bool("project", false, "")
var FlagRender = flag.Bool("render", false, "")
var FlagUtilities = flag.Bool("utilities", false, "")
var FlagViews = flag.String("views", "", "")
var FlagOut = flag.String("out", "", "")
