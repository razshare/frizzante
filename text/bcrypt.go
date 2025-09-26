package text

import "golang.org/x/crypto/bcrypt"

func BcryptHash(text string, cost int) (string, error) {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(text), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func BcryptHashDefault(text string) (string, error) {
	return BcryptHash(text, bcrypt.DefaultCost)
}

func BcryptCompare(hash string, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}
