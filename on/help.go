package on

import flag "github.com/spf13/pflag"

func Help() error {
	flag.Usage()
	return nil
}
