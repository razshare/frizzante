package on

import (
	"github.com/razshare/frizzante/cli"
	flag "github.com/spf13/pflag"
)

func Help(_ *cli.Cli) error {
	flag.Usage()
	return nil
}
