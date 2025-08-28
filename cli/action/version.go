package action

import "strings"

func Version(opts VersionOptions) (err error) {
	var data []byte
	if data, err = opts.Efs.ReadFile("version"); err != nil {
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
