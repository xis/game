package domain

import "context"

//go:generate mockery --name TokenManager --structname MockTokenManager --outpkg mocks --filename token_manager_mock.go --output ./mocks/. --with-expecter
type TokenManager interface {
	Create(ctx context.Context, userID string) (string, error)
	ExtractUserID(ctx context.Context, token string) (string, error)
}
