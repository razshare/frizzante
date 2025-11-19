package actions

import (
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/inputs"
)

func CreateProject(options CreateProjectOptions) (name string, err error) {
	name = options.Name
	if options.Name == "" {
		options.Name, err = inputs.Send("give the project a name")
		name = options.Name
		if err != nil {
			return
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
