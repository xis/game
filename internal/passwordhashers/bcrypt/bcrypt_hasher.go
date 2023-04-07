package bcrypt

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

func (hasher *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (hasher *BcryptPasswordHasher) ComparePasswordAndHash(password, hash string) (bool, error) {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}
