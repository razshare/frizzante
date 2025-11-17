package actions

import (
	"github.com/razshare/frizzante/cli/generate"
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

	return
}
