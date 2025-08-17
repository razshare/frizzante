package wrap

import "strings"

func Send(text string, width int) []string {
	if width <= 0 {
		return strings.Split(text, "\n")
	}

	inputLines := strings.Split(text, "\n")
	var result []string

	for _, line := range inputLines {
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

		currentLine := ""
		for _, word := range words {
			if currentLine == "" {
				currentLine = word
			} else if len(currentLine)+1+len(word) <= width {
				currentLine += " " + word
			} else {
				result = append(result, currentLine)
				currentLine = word
			}
		}
		if currentLine != "" {
			result = append(result, currentLine)
		}
	}

	return result
}
