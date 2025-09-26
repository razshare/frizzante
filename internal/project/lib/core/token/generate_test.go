package token

import (
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		userId   string
		hashType HashType
		secret   string
		wantErr  bool
	}{
		{
			name:     "SHA256 hash",
			userId:   "user123",
			hashType: SHA256,
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "SHA512 hash",
			userId:   "user456",
			hashType: SHA512,
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "SHA3_256 hash",
			userId:   "user789",
			hashType: SHA3_256,
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "SHA3_512 hash",
			userId:   "userABC",
			hashType: SHA3_512,
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "HMAC_SHA256 with secret",
			userId:   "userDEF",
			hashType: HMAC_SHA256,
			secret:   "my-secret-key",
			wantErr:  false,
		},
		{
			name:     "HMAC_SHA512 with secret",
			userId:   "userGHI",
			hashType: HMAC_SHA512,
			secret:   "another-secret-key",
			wantErr:  false,
		},
		{
			name:     "Invalid hash type defaults to SHA256",
			userId:   "userXYZ",
			hashType: HashType(999),
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "Empty userId",
			userId:   "",
			hashType: SHA256,
			secret:   "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := Generate(tt.userId, tt.hashType, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if token == "" {
					t.Error("Generate() returned empty token")
				}
				// Verify token is valid hex
				if _, err := hex.DecodeString(token); err != nil {
					t.Errorf("Generate() returned invalid hex: %v", err)
				}
			}
		})
	}
}

func TestGenerateDeterminism(t *testing.T) {
	// Test that same input produces different outputs due to randomness
	userId := "testuser"
	hashType := SHA256
	secret := ""

	token1, err1 := Generate(userId, hashType, secret)
	if err1 != nil {
		t.Fatalf("First Generate() failed: %v", err1)
	}

	token2, err2 := Generate(userId, hashType, secret)
	if err2 != nil {
		t.Fatalf("Second Generate() failed: %v", err2)
	}

	if token1 == token2 {
		t.Error("Generate() should produce different tokens due to randomness and timestamp")
	}
}

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		cost     int
		wantErr  bool
	}{
		{
			name:     "Normal password with default cost",
			password: "myPassword123",
			cost:     bcrypt.DefaultCost,
			wantErr:  false,
		},
		{
			name:     "Password with minimum cost",
			password: "anotherPass456",
			cost:     bcrypt.MinCost,
			wantErr:  false,
		},
		{
			name:     "Password with below minimum cost (should use default)",
			password: "testPass789",
			cost:     bcrypt.MinCost - 1,
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			cost:     bcrypt.DefaultCost,
			wantErr:  false,
		},
		{
			name:     "Very long password (exceeds bcrypt limit)",
			password: strings.Repeat("a", 100),
			cost:     bcrypt.DefaultCost,
			wantErr:  true,
		},
		{
			name:     "Password with special characters",
			password: "P@ssw0rd!#$%^&*()",
			cost:     bcrypt.DefaultCost,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GeneratePassword(tt.password, tt.cost)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("GeneratePassword() returned empty hash")
				}

				if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
					t.Error("GeneratePassword() returned invalid bcrypt hash format")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wrongPass string
		cost      int
	}{
		{
			name:      "Correct password verification",
			password:  "correctPassword123",
			wrongPass: "wrongPassword456",
			cost:      bcrypt.DefaultCost,
		},
		{
			name:      "Case sensitive verification",
			password:  "CaseSensitive",
			wrongPass: "casesensitive",
			cost:      bcrypt.MinCost,
		},
		{
			name:      "Similar but different passwords",
			password:  "password123",
			wrongPass: "password124",
			cost:      bcrypt.DefaultCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GeneratePassword(tt.password, tt.cost)
			if err != nil {
				t.Fatalf("GeneratePassword() failed: %v", err)
			}

			if !VerifyPassword(hash, tt.password) {
				t.Error("VerifyPassword() failed to verify correct password")
			}

			if VerifyPassword(hash, tt.wrongPass) {
				t.Error("VerifyPassword() incorrectly verified wrong password")
			}

			if VerifyPassword(hash, "") && tt.password != "" {
				t.Error("VerifyPassword() incorrectly verified empty password")
			}

			if VerifyPassword("invalid-hash", tt.password) {
				t.Error("VerifyPassword() incorrectly verified with invalid hash")
			}
		})
	}
}

func TestGenerateSecure(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Normal length",
			length:  32,
			wantErr: false,
		},
		{
			name:    "Small length",
			length:  1,
			wantErr: false,
		},
		{
			name:    "Large length",
			length:  256,
			wantErr: false,
		},
		{
			name:    "Zero length",
			length:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateSecure(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSecure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.length > 0 {
				if token == "" {
					t.Error("GenerateSecure() returned empty token")
				}

				if !IsValidBase64(token) {
					t.Error("GenerateSecure() returned invalid base64")
				}
			}
		})
	}
}

func TestGenerateSecureUniqueness(t *testing.T) {
	length := 32
	tokens := make(map[string]bool)
	iterations := 100

	for range iterations {
		token, err := GenerateSecure(length)
		if err != nil {
			t.Fatalf("GenerateSecure() failed: %v", err)
		}
		if tokens[token] {
			t.Error("GenerateSecure() produced duplicate token")
		}
		tokens[token] = true
	}
}

func TestEncodeToHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "Single byte",
			input:    []byte{0x0F},
			expected: "0f",
		},
		{
			name:     "Multiple bytes",
			input:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
			expected: "deadbeef",
		},
		{
			name:     "All zeros",
			input:    []byte{0x00, 0x00, 0x00},
			expected: "000000",
		},
		{
			name:     "All ones",
			input:    []byte{0xFF, 0xFF},
			expected: "ffff",
		},
		{
			name:     "Mixed values",
			input:    []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0},
			expected: "123456789abcdef0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeToHex(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeToHex() = %v, want %v", result, tt.expected)
			}
			// Verify it's valid hex
			decoded, err := hex.DecodeString(result)
			if err != nil {
				t.Errorf("EncodeToHex() produced invalid hex: %v", err)
			}
			// Verify round-trip
			if string(decoded) != string(tt.input) {
				t.Error("EncodeToHex() round-trip failed")
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	hashTypes := []struct {
		name     string
		hashType HashType
		secret   string
	}{
		{"SHA256", SHA256, ""},
		{"SHA512", SHA512, ""},
		{"SHA3_256", SHA3_256, ""},
		{"SHA3_512", SHA3_512, ""},
		{"HMAC_SHA256", HMAC_SHA256, "secret"},
		{"HMAC_SHA512", HMAC_SHA512, "secret"},
	}

	for _, ht := range hashTypes {
		b.Run(ht.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = Generate("user123", ht.hashType, ht.secret)
			}
		})
	}
}

func BenchmarkGeneratePassword(b *testing.B) {
	costs := []int{bcrypt.MinCost, bcrypt.DefaultCost}
	password := "testPassword123"

	for _, cost := range costs {
		b.Run("Cost"+string(rune(cost)), func(b *testing.B) {
			for b.Loop() {
				_, _ = GeneratePassword(password, cost)
			}
		})
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testPassword123"
	hash, _ := GeneratePassword(password, bcrypt.DefaultCost)

	b.ResetTimer()
	for b.Loop() {
		_ = VerifyPassword(hash, password)
	}
}

func BenchmarkEncodeToHex(b *testing.B) {
	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for b.Loop() {
		_ = EncodeToHex(data)
	}
}

func IsValidBase64(s string) bool {
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	for _, c := range s {
		if !strings.Contains(validChars, string(c)) {
			return false
		}
	}
	return true
}
