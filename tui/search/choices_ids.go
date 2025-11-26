package search

func ChoicesIds(choices []Choice) []string {
	result := make([]string, len(choices))
	for index, choice := range choices {
		result[index] = choice.Id
	}
	return result
}
