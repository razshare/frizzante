package actions

import flag "github.com/spf13/pflag"

func Help(_ HelpOptions) (err error) {
	flag.Usage()
	return nil
}
