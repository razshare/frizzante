package actions

import "github.com/razshare/frizzante/tui/messages"

func Configure(options ConfigureOptions) (err error) {
	if err = Install(InstallOptions{
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}

	if err = Package(PackageOptions{
		Bun:        options.Bun,
		Production: false,
	}); err != nil {
		return
	}

	messages.Success("project configured")

	return
}
