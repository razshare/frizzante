package actions

import (
	"os"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/platforms"
	"github.com/razshare/frizzante/tui/inputs"
)

func CreateProject(options CreateProjectOptions) (err error) {
	if options.Name == "" {
		options.Name, err = inputs.Send("give the project a name")
		if err != nil {
			return err
		}
	}

	if err = generate.Project(generate.ProjectOptions{
		Name: options.Name,
		Go:   options.Go,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = os.Chdir(options.Name); err != nil {
		return
	}

	if err = Configure(ConfigureOptions{
		Auto:     true,
		Go:       options.Go,
		Efs:      options.Efs,
		Platform: platforms.Detect(),
		Bun:      options.Bun,
		Air:      options.Air,
	}); err != nil {
		return
	}

	if err = Install(InstallOptions{
		Go:  options.Go,
		Bun: options.Bun,
	}); err != nil {
		return
	}

	if err = Package(PackageOptions{
		Bun:  options.Bun,
		Prod: false,
	}); err != nil {
		return
	}

	//var frizzante string
	//if frizzante, err = os.Executable(); err != nil {
	//	return
	//}
	//
	//if !messages.Command("", os.Environ(), frizzante, "--configure") {
	//	err = errors.New("could not configure project")
	//	return
	//}
	//
	//if !messages.Command("", os.Environ(), frizzante, "--install") {
	//	err = errors.New("could not install packages")
	//	return
	//}
	//
	//if !messages.Command("", os.Environ(), frizzante, "--package") {
	//	err = errors.New("could not package app")
	//	return
	//}

	return
}
