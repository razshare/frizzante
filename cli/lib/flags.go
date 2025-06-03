package lib

import "flag"

var FlagGenerate = flag.Bool("generate", true, "")
var FlagProject = flag.Bool("project", false, "")
var FlagRouter = flag.Bool("router", false, "")
var FlagUtilities = flag.Bool("utilities", false, "")
var FlagViews = flag.String("views", "", "")
var FlagOut = flag.String("out", "", "")
