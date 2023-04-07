package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"game/internal/domain"
	"game/internal/domain/mocks"
)

type UserServiceTestSuite struct {
	suite.Suite

	service *userService

	mockUserScoreRepository *mocks.MockUserScoreRepository
	mockUserRepository      *mocks.MockUserRepository
	mockTokenManager        *mocks.MockTokenManager
	mockPasswordHasher      *mocks.MockPasswordHasher
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.mockUserScoreRepository = mocks.NewMockUserScoreRepository(suite.T())
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.T())
	suite.mockTokenManager = mocks.NewMockTokenManager(suite.T())
	suite.mockPasswordHasher = mocks.NewMockPasswordHasher(suite.T())

	suite.service = NewUserService(UserServiceDependencies{
		UserRepository:      suite.mockUserRepository,
		UserScoreRepository: suite.mockUserScoreRepository,
		TokenManager:        suite.mockTokenManager,
		PasswordHasher:      suite.mockPasswordHasher,
	})
}

func (suite *UserServiceTestSuite) TestLogin() {
	suite.mockUserRepository.
		EXPECT().
		GetByName(mock.Anything, "username").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password-hash",
		}, nil)

	suite.mockPasswordHasher.
		EXPECT().
		ComparePasswordAndHash("password", "password-hash").
		Return(true, nil)

	suite.mockTokenManager.
		EXPECT().
		Create(mock.Anything, "user-id").
		Return("token", nil)

	result, err := suite.service.Login(context.Background(), "username", "password")
	suite.NoError(err)

	suite.Equal("token", result.Token)
	suite.Equal("user-id", result.UserID)
	suite.Equal("username", result.UserName)
}

func (suite *UserServiceTestSuite) TestLogin_InvalidPassword() {
	suite.mockUserRepository.
		EXPECT().
		GetByName(mock.Anything, "username").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password-hash",
		}, nil)

	suite.mockPasswordHasher.
		EXPECT().
		ComparePasswordAndHash("password", "password-hash").
		Return(false, nil)

	result, err := suite.service.Login(context.Background(), "username", "password")
	suite.Equal(ErrInvalidPassword, err)
	suite.Equal(LoginResult{}, result)
}

func (suite *UserServiceTestSuite) TestLogin_UserNotFound() {
	suite.mockUserRepository.
		EXPECT().
		GetByName(mock.Anything, "username").
		Return(domain.User{}, domain.ErrResourceNotFound)

	result, err := suite.service.Login(context.Background(), "username", "password")

	suite.Equal(err, domain.ErrResourceNotFound)
	suite.Empty(result)
}

func (suite *UserServiceTestSuite) TestLogin_PasswordHasherFailed() {
	suite.mockUserRepository.
		EXPECT().
		GetByName(mock.Anything, "username").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password-hash",
		}, nil)

	suite.mockPasswordHasher.
		EXPECT().
		ComparePasswordAndHash("password", "password-hash").
		Return(false, domain.ErrInternal)

	result, err := suite.service.Login(context.Background(), "username", "password")

	suite.Equal(err, domain.ErrInternal)
	suite.Empty(result)
}

func (suite *UserServiceTestSuite) TestLogin_TokenManagerFailed() {
	suite.mockUserRepository.
		EXPECT().
		GetByName(mock.Anything, "username").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password-hash",
		}, nil)

	suite.mockPasswordHasher.
		EXPECT().
		ComparePasswordAndHash("password", "password-hash").
		Return(true, nil)

	suite.mockTokenManager.
		EXPECT().
		Create(mock.Anything, "user-id").
		Return("", domain.ErrInternal)

	result, err := suite.service.Login(context.Background(), "username", "password")

	suite.Equal(err, domain.ErrInternal)
	suite.Empty(result)
}

func (suite *UserServiceTestSuite) TestRegister() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByName(mock.Anything, "username").
		Return(false, nil)

	suite.mockPasswordHasher.
		EXPECT().
		HashPassword("password").
		Return("password-hash", nil)

	suite.mockUserRepository.
		EXPECT().
		Create(mock.Anything, "username", "password-hash").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password-hash",
		}, nil)

	user, err := suite.service.Register(context.Background(), "username", "password")

	expectedResult := domain.User{
		ID:           "user-id",
		Name:         "username",
		PasswordHash: "password-hash",
	}

	suite.NoError(err)
	suite.Equal(expectedResult, user)
}

func (suite *UserServiceTestSuite) TestRegister_CheckExistsByNameFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByName(mock.Anything, "username").
		Return(false, domain.ErrInternal)

	user, err := suite.service.Register(context.Background(), "username", "password")

	suite.ErrorIs(err, domain.ErrInternal)
	suite.Empty(user)
}

func (suite *UserServiceTestSuite) TestRegister_UsernameExists() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByName(mock.Anything, "username").
		Return(true, nil)

	user, err := suite.service.Register(context.Background(), "username", "password")

	suite.ErrorIs(err, ErrUsernameExists)
	suite.Empty(user)
}

func (suite *UserServiceTestSuite) TestRegister_PasswordHasherFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByName(mock.Anything, "username").
		Return(false, nil)

	suite.mockPasswordHasher.
		EXPECT().
		HashPassword("password").
		Return("", domain.ErrInternal)

	user, err := suite.service.Register(context.Background(), "username", "password")

	suite.ErrorIs(err, domain.ErrInternal)
	suite.Empty(user)
}

func (suite *UserServiceTestSuite) TestRegister_UserRepositoryFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByName(mock.Anything, "username").
		Return(false, nil)

	suite.mockPasswordHasher.
		EXPECT().
		HashPassword("password").
		Return("password-hash", nil)

	suite.mockUserRepository.
		EXPECT().
		Create(mock.Anything, "username", "password-hash").
		Return(domain.User{}, domain.ErrInternal)

	user, err := suite.service.Register(context.Background(), "username", "password")

	suite.ErrorIs(err, domain.ErrInternal)
	suite.Empty(user)
}
