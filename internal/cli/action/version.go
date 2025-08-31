package action

import "strings"

func Version(options VersionOptions) (err error) {
	var data []byte
	if data, err = options.Efs.ReadFile("version"); err != nil {
		return
	}

	version := string(data)

	lines := strings.Split(version, "\n")

	if len(lines) == 0 {
		return
	}

	println(lines[0])

	return
}
