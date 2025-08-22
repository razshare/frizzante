package action

import "strings"

func Version(opts VersionOptions) error {
	data, err := opts.Efs.ReadFile("version")
	if err != nil {
		return err
	}

	version := string(data)

	lines := strings.Split(version, "\n")

	if len(lines) == 0 {
		return nil
	}

	println(lines[0])

	return nil
}
