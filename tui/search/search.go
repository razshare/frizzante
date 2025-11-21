package search

type Search struct {
	Value    string
	Choices  []Choice
	Filtered []Choice
	Active   bool
}
