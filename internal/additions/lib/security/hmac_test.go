package security

import "testing"

func TestHmacSha256(t *testing.T) {
	signature := HmacSha256("hello", "secretkey")

	if signature != "122b99e68dd9cabbd464c943550399cad150790bcd3d94f526b92fa29fb762bc" {
		t.Fatal("hmac sha256 signature should be 122b99e68dd9cabbd464c943550399cad150790bcd3d94f526b92fa29fb762bc")
	}

	if !HmacVerifySha256("hello", "secretkey", signature) {
		t.Fatal("hmac sha256 verification should succeed")
	}

	if HmacVerifySha256("hello", "wrongkey", signature) {
		t.Fatal("hmac sha256 verification should fail with wrong key")
	}
}

func TestHmacSha512(t *testing.T) {
	signature := HmacSha512("hello", "secretkey")

	if signature != "2fcb078bae0caf610e46efb16216693ce8bb445e7965d66db175ba66356e3da7a5b4032ecd0b3c31b6b71c2c199cef50e06824649c126766e21fd2099a0cf438" {
		t.Fatal("hmac sha256 signature should be 2fcb078bae0caf610e46efb16216693ce8bb445e7965d66db175ba66356e3da7a5b4032ecd0b3c31b6b71c2c199cef50e06824649c126766e21fd2099a0cf438")
	}

	if !HmacVerifySha512("hello", "secretkey", signature) {
		t.Fatal("hmac sha256 verification should succeed")
	}

	if HmacVerifySha512("hello", "wrongkey", signature) {
		t.Fatal("hmac sha256 verification should fail with wrong key")
	}
}
