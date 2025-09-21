package text

import "crypto/sha1"

func Sha1(text string) string {
	from := sha1.Sum([]byte(text))
	to := make([]byte, 40)
	j := 0
	for _, v := range from {
		to[j] = HexTable[v>>4]
		to[j+1] = HexTable[v&0x0f]
		j += 2
	}
	return string(to)
}
