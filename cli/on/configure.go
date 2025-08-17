package on

func Configure(app string, clr bool, yes bool, gobin string, bunbin string, sqlcbin string) error {
	err := Generate(app, "bun,air", clr, yes, gobin, sqlcbin)
	if err != nil {
		return err
	}
	return Install(app, gobin, bunbin)
}
