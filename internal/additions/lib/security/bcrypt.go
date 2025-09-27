package security

import "golang.org/x/crypto/bcrypt"

// BcryptHash returns the bcrypt hash of the password with a default cost of 10.
func BcryptHash(text string) (string, error) {
	return BcryptHashWithCost(text, bcrypt.DefaultCost)
}

// BcryptHashWithCost returns the bcrypt hash of the password at the given cost.
func BcryptHashWithCost(text string, cost int) (string, error) {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(text), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BcryptVerify verifies that a bcrypt hash matches the given text.
func BcryptVerify(hash string, text string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(text)) == nil
}
