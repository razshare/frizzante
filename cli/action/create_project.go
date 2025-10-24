package action

import (
	"errors"
	"os"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
)

func CreateProject(options CreateProjectOptions) (err error) {
	if options.Name == "" {
		options.Name, err = input.Send("give the project a name")
		if err != nil {
			return err
		}
	}

	err = generate.Project(generate.ProjectOptions{
		Name: options.Name,
		Go:   options.Go,
		Efs:  options.Efs,
	})

	var frizzante string
	if frizzante, err = os.Executable(); err != nil {
		return
	}

	if !messages.Command(options.Name, os.Environ(), frizzante, "--configure") {
		err = errors.New("could not configure project")
		return
	}

	if !messages.Command(options.Name, os.Environ(), frizzante, "--install") {
		err = errors.New("could not install packages")
		return
	}

	if !messages.Command(options.Name, os.Environ(), frizzante, "--package") {
		err = errors.New("could not package app")
		return
	}

	return
}
