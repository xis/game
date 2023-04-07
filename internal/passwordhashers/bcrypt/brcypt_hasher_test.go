package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordHasher(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hash, err := hasher.HashPassword("password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	match, err := hasher.ComparePasswordAndHash("password", hash)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestBcryptPasswordHasher_InvalidPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	hash, err := hasher.HashPassword("password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	match, err := hasher.ComparePasswordAndHash("invalid-password", hash)
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestBcryptPasswordHasher_LongPassword(t *testing.T) {
	hasher := NewBcryptPasswordHasher()

	password := ""

	for i := 0; i < 100; i++ {
		password += "password"
	}

	hash, err := hasher.HashPassword(password)
	assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	assert.Empty(t, hash)
}
