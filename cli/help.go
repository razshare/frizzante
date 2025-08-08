package cli

import flag "github.com/spf13/pflag"

func OnHelp() {
	flag.Usage()
}
