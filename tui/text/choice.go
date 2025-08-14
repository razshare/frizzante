package text

import "strings"

// TitleAndContent separates the first line from the rest of the text.
//
// Lines are trimmed.
func TitleAndContent(txt string) (tlt string, cnt string) {
	p := strings.SplitN(strings.TrimSpace(txt), "\n", 2)
	if len(p) > 1 {
		for _, l := range strings.Split(p[1], "\n") {
			trm := strings.TrimSpace(l)
			if trm == "" {
				continue
			}
			cnt = cnt + trm + " "
		}
	} else {
		cnt = ""
	}
	tlt = p[0]
	return
}
