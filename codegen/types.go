package codegen

type State uint64

const Start State = 0
const EscapingOriginal State = 98
const EscapingReplacement State = 99
const ReadingOriginalString State = 100
const DoneReadingOriginalString State = 101
const ReadingReplacementString State = 200
const DoneReadingReplacementString State = 201
const ExpectingReplacementString State = 900
const Invalid State = 1000

type Mod struct {
	Pattern     string
	Replacement string
}

type Submit func(c string)
type Build func(s Section) error
type Section struct {
	Mods []Mod
	Line *string
}

type Generation struct {
	From      string
	To        string
	Overwrite func(n string) bool
}
