package text

import "crypto/sha1"

func Sha1(text string) string {
	from := sha1.Sum([]byte(text))
	to := make([]byte, 40) // must be 40, double the size of sha1 sum (20)
	var i int
	for _, v := range from {
		to[i] = HexTable[v>>4]
		to[i+1] = HexTable[v&0x0f]
		i += 2
	}
	return string(to)
}
