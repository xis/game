package grpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"game/internal/domain"
	userpb "game/internal/proto/user/proto"
	"game/internal/services"
	"game/internal/services/mocks"
)

type UserControllerTestSuite struct {
	suite.Suite

	controller *userController

	mockUserService *mocks.MockUserService
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.mockUserService = mocks.NewMockUserService(suite.T())

	suite.controller = NewUserController(UserControllerDependencies{
		UserService: suite.mockUserService,
		Logger:      logrus.New(),
	})
}

func (suite *UserControllerTestSuite) TestLogin() {
	suite.mockUserService.
		EXPECT().
		Login(mock.Anything, "username", "password").
		Return(services.LoginResult{
			UserID:   "user-id",
			UserName: "username",
			Token:    "token",
		}, nil)

	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "username",
		Password: "password",
	})
	suite.NoError(err)

	expectedResult := &userpb.LoginResponse{
		Status: StatusSuccess,
		Result: &userpb.LoginResult{
			Token:    "token",
			Username: "username",
			UserID:   "user-id",
		},
	}

	suite.Equal(expectedResult.Status, result.Status)
	suite.Equal(expectedResult.Result, result.Result)
	suite.NotEmpty(result.Timestamp)
}

func (suite *UserControllerTestSuite) TestLogin_NoUsername() {
	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "",
		Password: "password",
	})

	suite.ErrorIs(err, ErrUsernameRequired)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestLogin_NoPassword() {
	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "username",
		Password: "",
	})

	suite.ErrorIs(err, ErrPasswordRequired)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestLogin_InvalidPassword() {
	suite.mockUserService.
		EXPECT().
		Login(mock.Anything, "username", "password").
		Return(services.LoginResult{}, services.ErrInvalidPassword)

	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "username",
		Password: "password",
	})

	suite.ErrorIs(err, ErrInvalidPassword)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestLogin_UserNotFound() {
	suite.mockUserService.
		EXPECT().
		Login(mock.Anything, "username", "password").
		Return(services.LoginResult{}, domain.ErrResourceNotFound)

	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "username",
		Password: "password",
	})

	suite.ErrorIs(err, ErrUserNotFound)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestLogin_InternalError() {
	suite.mockUserService.
		EXPECT().
		Login(mock.Anything, "username", "password").
		Return(services.LoginResult{}, domain.ErrInternal)

	result, err := suite.controller.Login(context.Background(), &userpb.LoginRequest{
		Username: "username",
		Password: "password",
	})

	suite.ErrorIs(err, ErrInternal)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestRegister() {
	suite.mockUserService.
		EXPECT().
		Register(mock.Anything, "username", "password").
		Return(domain.User{
			ID:           "user-id",
			Name:         "username",
			PasswordHash: "password",
		}, nil)

	result, err := suite.controller.Register(context.Background(), &userpb.RegisterRequest{
		Username: "username",
		Password: "password",
	})
	suite.NoError(err)

	expectedResult := &userpb.RegisterResponse{
		Status: StatusSuccess,
		Result: &userpb.RegistrationResult{
			Username: "username",
			Password: "password",
			UserID:   "user-id",
		},
	}

	suite.Equal(expectedResult.Status, result.Status)
	suite.Equal(expectedResult.Result, result.Result)
	suite.NotEmpty(result.Timestamp)
}

func (suite *UserControllerTestSuite) TestRegister_NoUsername() {
	result, err := suite.controller.Register(context.Background(), &userpb.RegisterRequest{
		Username: "",
		Password: "password",
	})

	suite.ErrorIs(err, ErrUsernameRequired)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestRegister_NoPassword() {
	result, err := suite.controller.Register(context.Background(), &userpb.RegisterRequest{
		Username: "username",
		Password: "",
	})

	suite.ErrorIs(err, ErrPasswordRequired)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestRegister_UsernameAlreadyExists() {
	suite.mockUserService.
		EXPECT().
		Register(mock.Anything, "username", "password").
		Return(domain.User{}, services.ErrUsernameExists)

	result, err := suite.controller.Register(context.Background(), &userpb.RegisterRequest{
		Username: "username",
		Password: "password",
	})

	suite.ErrorIs(err, ErrUsernameExists)
	suite.Empty(result)
}

func (suite *UserControllerTestSuite) TestRegister_InternalError() {
	suite.mockUserService.
		EXPECT().
		Register(mock.Anything, "username", "password").
		Return(domain.User{}, domain.ErrInternal)

	result, err := suite.controller.Register(context.Background(), &userpb.RegisterRequest{
		Username: "username",
		Password: "password",
	})

	suite.ErrorIs(err, ErrInternal)
	suite.Empty(result)
}
