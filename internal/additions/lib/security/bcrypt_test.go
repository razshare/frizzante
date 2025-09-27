package security

import "testing"

func TestBcryptHash(t *testing.T) {
	password := "hello"

	hash, err := BcryptHash(password)
	if err != nil {
		t.Fatal("error generating bcrypt hash:", err)
	}

	if !BcryptVerify(hash, password) {
		t.Fatal("bcrypt compare should return true for correct password")
	}

	if BcryptVerify(hash, "wrong") {
		t.Fatal("bcrypt compare should return false for wrong password")
	}
}

func TestBcryptHashWithCost(t *testing.T) {
	password := "hello"

	hash, err := BcryptHashWithCost(password, 4)
	if err != nil {
		t.Fatal("error generating bcrypt hash with custom cost:", err)
	}

	if !BcryptVerify(hash, password) {
		t.Fatal("bcrypt compare should return true for correct password with custom cost")
	}
}
