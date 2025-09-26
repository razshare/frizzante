package text

import "testing"

func TestBcrypt(t *testing.T) {
	password := "hello"

	hash, err := BcryptHashDefault(password)
	if err != nil {
		t.Fatal("error generating bcrypt hash:", err)
	}

	if !BcryptCompare(hash, password) {
		t.Fatal("bcrypt compare should return true for correct password")
	}

	if BcryptCompare(hash, "wrong") {
		t.Fatal("bcrypt compare should return false for wrong password")
	}
}

func TestBcryptWithCustomCost(t *testing.T) {
	password := "hello"

	hash, err := BcryptHash(password, 4)
	if err != nil {
		t.Fatal("error generating bcrypt hash with custom cost:", err)
	}

	if !BcryptCompare(hash, password) {
		t.Fatal("bcrypt compare should return true for correct password with custom cost")
	}
}
