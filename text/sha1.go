package text

import "crypto/sha1"

func Sha1(text string) string {
	from := sha1.Sum([]byte(text))
	to := make([]byte, 40) // must be 40, double the size of sha1 sum (20)
	var i int
	for _, v := range from {
		to[i] = HexTable[v>>4]     // encoding the left half, first nibble
		to[i+1] = HexTable[v&0x0f] // encoding the right half, second nibble
		i += 2                     // jump by 2 because we encoded 2 nibbles
	}
	return string(to)
}
