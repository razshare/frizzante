package action

import "strings"

func Version(o VersionOptions) error {
	d, err := o.Efs.ReadFile("version")
	if err != nil {
		return err
	}

	v := string(d)

	ls := strings.Split(v, "\n")

	if len(ls) == 0 {
		return nil
	}

	println(ls[0])

	return nil
}
