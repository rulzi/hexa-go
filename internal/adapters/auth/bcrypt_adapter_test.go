package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBcryptPasswordHasher(t *testing.T) {
	hasher := NewBcryptPasswordHasher()
	assert.NotNil(t, hasher)
}

func TestBcryptPasswordHasher_Hash_Success(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := "testpassword123"
	hashed, err := hasher.Hash(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed) // Hashed should be different from original
	assert.Len(t, hashed, 60)            // bcrypt hash is always 60 characters
}

func TestBcryptPasswordHasher_Hash_DifferentPasswords(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password1 := "password1"
	password2 := "password2"

	hashed1, err1 := hasher.Hash(password1)
	hashed2, err2 := hasher.Hash(password2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hashed1, hashed2) // Different passwords should produce different hashes
}

func TestBcryptPasswordHasher_Hash_SamePasswordDifferentHashes(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := "samepassword"
	hashed1, err1 := hasher.Hash(password)
	hashed2, err2 := hasher.Hash(password)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	// Same password should produce different hashes due to salt
	assert.NotEqual(t, hashed1, hashed2)
}

func TestBcryptPasswordHasher_Hash_EmptyPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hashed, err := hasher.Hash("")

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
}

func TestBcryptPasswordHasher_Hash_LongPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	// Create a long password (but within bcrypt's 72 byte limit)
	// Using 70 characters to be safe
	longPassword := ""
	for i := 0; i < 70; i++ {
		longPassword += "a"
	}

	hashed, err := hasher.Hash(longPassword)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
}

func TestBcryptPasswordHasher_Hash_PasswordExceedsLimit(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	// Create a password that exceeds bcrypt's 72 byte limit
	longPassword := ""
	for i := 0; i < 100; i++ {
		longPassword += "a"
	}

	hashed, err := hasher.Hash(longPassword)

	// Bcrypt should return an error for passwords exceeding 72 bytes
	assert.Error(t, err)
	assert.Empty(t, hashed)
	assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
}

func TestBcryptPasswordHasher_Verify_Success(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := "testpassword123"
	hashed, err := hasher.Hash(password)
	assert.NoError(t, err)

	valid := hasher.Verify(hashed, password)
	assert.True(t, valid)
}

func TestBcryptPasswordHasher_Verify_WrongPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := "testpassword123"
	wrongPassword := "wrongpassword"
	hashed, err := hasher.Hash(password)
	assert.NoError(t, err)

	valid := hasher.Verify(hashed, wrongPassword)
	assert.False(t, valid)
}

func TestBcryptPasswordHasher_Verify_EmptyPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := "testpassword"
	hashed, err := hasher.Hash(password)
	assert.NoError(t, err)

	valid := hasher.Verify(hashed, "")
	assert.False(t, valid)
}

func TestBcryptPasswordHasher_Verify_InvalidHash(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	invalidHash := "invalidhash"
	password := "testpassword"

	valid := hasher.Verify(invalidHash, password)
	assert.False(t, valid)
}

func TestBcryptPasswordHasher_Verify_EmptyHash(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	valid := hasher.Verify("", "testpassword")
	assert.False(t, valid)
}

func TestBcryptPasswordHasher_HashAndVerify_RoundTrip(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	testCases := []string{
		"simple",
		"password123",
		"P@ssw0rd!",
		"very-long-password-with-special-chars-!@#$%^&*()",
		"1234567890",
		"   spaces   ",
		"unicode-测试-🚀",
	}

	for _, password := range testCases {
		t.Run(password, func(t *testing.T) {
			hashed, err := hasher.Hash(password)
			assert.NoError(t, err)

			valid := hasher.Verify(hashed, password)
			assert.True(t, valid, "Password %s should verify correctly", password)
		})
	}
}

