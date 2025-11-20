package actions

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
)

func CreateProject(options CreateProjectOptions) (err error) {
	var name string
	if name = options.Name; name == "" {
		name, err = inputs.Send("give the project a name")
		if err != nil {
			return
		}
	}

	if err = generate.Project(generate.ProjectOptions{
		Name: name,
		Go:   options.Go,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	messages.Successf("project %s created with success!", name)

	step1 := fmt.Sprintf("1. cd %s", name)
	step2 := fmt.Sprintf("2. frizzante --configure")
	step3 := fmt.Sprintf("3. frizzante --dev")

	length1 := len(step1)
	length2 := len(step2)
	length3 := len(step3)

	width := length2

	if length1 > length2 {
		width = length1
	}

	padding1 := strings.Repeat(" ", width-length1)
	padding2 := strings.Repeat(" ", width-length2)
	padding3 := strings.Repeat(" ", width-length3)

	comment1 := fmt.Sprintf("%s#changes directory to %s", padding1, name)
	comment2 := fmt.Sprintf("%s#configures the project", padding2)
	comment3 := fmt.Sprintf("%s#starts development mode", padding3)

	messages.Tipf(strings.Join([]string{
		"## next steps",
		fmt.Sprintf("%s %s", step1, comment1),
		fmt.Sprintf("%s %s", step2, comment2),
		fmt.Sprintf("%s %s", step3, comment3),
	}, "\n"))

	return
}
