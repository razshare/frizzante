package wrap

import "strings"

func Send(text string, width int) []string {
	if width <= 0 {
		return strings.Split(text, "\n")
	}

	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		if line == "" {
			continue
		}

		if len(line) <= width {
			result = append(result, line)
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			result = append(result, line)
			continue
		}

		cline := ""
		for _, word := range words {
			if cline == "" {
				cline = word
			} else if len(cline)+1+len(word) <= width {
				cline += " " + word
			} else {
				result = append(result, cline)
				cline = word
			}
		}
		if cline != "" {
			result = append(result, cline)
		}
	}

	return result
}
