package cli

import "flag"

var FlagGenerate = flag.Bool("generate", true, "")
var FlagProject = flag.Bool("project", false, "")
var FlagUtilities = flag.Bool("utilities", false, "")
var FlagOut = flag.String("out", "", "")
