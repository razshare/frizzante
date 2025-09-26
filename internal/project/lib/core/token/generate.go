package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"strconv"
	"time"

	"github.com/razshare/frizzante/text"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

const (
	SHA256 HashType = iota
	SHA512
	SHA3_256
	SHA3_512
	HMAC_SHA256
	HMAC_SHA512
)

func Generate(userId string, hashType HashType, secret string) (token string, err error) {
	randomBytes := make([]byte, 32)
	if _, err = rand.Read(randomBytes); err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	dataLen := len(userId) + 1 + len(timestamp) + 1 + (len(randomBytes) * 2)
	data := make([]byte, 0, dataLen)

	data = append(data, userId...)
	data = append(data, ':')

	data = append(data, timestamp...)
	data = append(data, ':')

	for _, b := range randomBytes {
		data = append(data, text.HexTable[b>>4])
		data = append(data, text.HexTable[b&0x0f])
	}

	switch hashType {
	case SHA256:
		hash := sha256.Sum256(data)
		return EncodeToHex(hash[:]), nil

	case SHA512:
		hash := sha512.Sum512(data)
		return EncodeToHex(hash[:]), nil

	case SHA3_256:
		hash := sha3.Sum256(data)
		return EncodeToHex(hash[:]), nil

	case SHA3_512:
		hash := sha3.Sum512(data)
		return EncodeToHex(hash[:]), nil

	case HMAC_SHA256:
		h := hmac.New(sha256.New, []byte(secret))
		h.Write(data)
		return EncodeToHex(h.Sum(nil)), nil

	case HMAC_SHA512:
		h := hmac.New(sha512.New, []byte(secret))
		h.Write(data)
		return EncodeToHex(h.Sum(nil)), nil

	default:
		return Generate(userId, SHA256, "")
	}
}

func GeneratePassword(password string, cost int) (hash string, err error) {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func VerifyPassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateSecure(length int) (token string, err error) {
	return text.RandomBase64(length)
}

func EncodeToHex(data []byte) string {
	result := make([]byte, len(data)*2)
	j := 0
	for _, b := range data {
		result[j] = text.HexTable[b>>4]
		result[j+1] = text.HexTable[b&0x0f]
		j += 2
	}
	return string(result)
}
