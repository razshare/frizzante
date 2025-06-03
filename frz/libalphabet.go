package frz

import "slices"

var alphabet = []rune{
	'A',
	'a',
	'B',
	'b',
	'C',
	'c',
	'D',
	'd',
	'E',
	'e',
	'F',
	'f',
	'G',
	'g',
	'H',
	'h',
	'I',
	'i',
	'J',
	'j',
	'K',
	'k',
	'L',
	'l',
	'M',
	'm',
	'N',
	'n',
	'O',
	'o',
	'P',
	'p',
	'Q',
	'q',
	'R',
	'r',
	'S',
	's',
	'T',
	't',
	'U',
	'u',
	'V',
	'v',
	'W',
	'w',
	'X',
	'x',
	'Y',
	'y',
	'Z',
	'z',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'0',
	'_',
	'-',
	'.',
}

// KeyIsSafe checks if a key is accepted by the safe alphabet.
//
// The safe alphabet accepts only "_" (underscore), "-" (hyphen), "." (period),
// english letters (uppercase and lowercase) and numeric digits (from 0 to 9).
func KeyIsSafe(key string) bool {
	for _, char := range key {
		if slices.Contains(alphabet, char) {
			continue
		}
		return false
	}

	return true
}
