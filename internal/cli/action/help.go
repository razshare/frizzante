package action

import flag "github.com/spf13/pflag"

func Help(_ HelpOptions) error {
	flag.Usage()
	return nil
}
