package domain

//go:generate mockery --name PasswordHasher --structname MockPasswordHasher --outpkg mocks --filename password_hasher_mock.go --output ./mocks/. --with-expecter
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswordAndHash(password, hash string) (bool, error)
}
