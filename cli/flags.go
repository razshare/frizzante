package cli

import "flag"

var FlagDevelop = flag.Bool("develop", false, "")
var FlagCreateProject = flag.Bool("project", false, "")
var FlagGenerateUtilities = flag.Bool("utilities", false, "")
var FlagOut = flag.String("out", "", "")
