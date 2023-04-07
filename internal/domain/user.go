package domain

import "context"

type User struct {
	ID           string
	Name         string
	PasswordHash string
}

//go:generate mockery --name UserRepository --structname MockUserRepository --outpkg mocks --filename user_repository_mock.go --output ./mocks/. --with-expecter
type UserRepository interface {
	Create(ctx context.Context, username, password string) (User, error)
	GetByName(ctx context.Context, username string) (User, error)
	CheckExistsByID(ctx context.Context, id string) (bool, error)
	CheckExistsByName(ctx context.Context, username string) (bool, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]User, error)
}
