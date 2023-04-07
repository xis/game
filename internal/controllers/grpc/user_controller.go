package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"game/internal/domain"
	userpb "game/internal/proto/user/proto"
	"game/internal/services"
)

const StatusSuccess = "success"

var (
	ErrInvalidPassword  = status.New(codes.InvalidArgument, "invalid password").Err()
	ErrInternal         = status.New(codes.Internal, "internal error").Err()
	ErrUsernameExists   = status.New(codes.AlreadyExists, "username exists").Err()
	ErrUsernameRequired = status.New(codes.InvalidArgument, "username is required").Err()
	ErrPasswordRequired = status.New(codes.InvalidArgument, "password is required").Err()
	ErrUserNotFound     = status.New(codes.NotFound, "user not found").Err()
)

type UserControllerDependencies struct {
	UserService services.UserService

	Logger *logrus.Logger
}

type userController struct {
	userpb.UnimplementedUserServiceServer

	userService services.UserService

	logger *logrus.Logger
}

func NewUserController(deps UserControllerDependencies) *userController {
	return &userController{
		userService: deps.UserService,
		logger:      deps.Logger,
	}
}

func (controller *userController) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	controller.
		logger.
		WithFields(logrus.Fields{
			"username": req.Username,
		}).
		Info("login request has been received")

	if req.Username == "" {
		return nil, ErrUsernameRequired
	}

	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	result, err := controller.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		controller.logger.
			WithFields(logrus.Fields{
				"username": req.Username,
			}).
			Error("login request is failed, ", err)

		if errors.Is(err, services.ErrInvalidPassword) {
			return nil, ErrInvalidPassword
		}

		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternal
	}

	return &userpb.LoginResponse{
		Status:    StatusSuccess,
		Timestamp: time.Now().Unix(),
		Result: &userpb.LoginResult{
			Token:    result.Token,
			Username: result.UserName,
			UserID:   result.UserID,
		},
	}, nil
}

func (controller *userController) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	controller.
		logger.
		WithFields(logrus.Fields{
			"username": req.Username,
		}).
		Info("register request has been received")

	if req.Username == "" {
		return nil, ErrUsernameRequired
	}

	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	userDetails, err := controller.userService.Register(ctx, req.Username, req.Password)
	if err != nil {
		controller.logger.
			WithFields(logrus.Fields{
				"username": req.Username,
			}).
			Error("register request is failed, ", err)

		if errors.Is(err, services.ErrUsernameExists) {
			return nil, ErrUsernameExists
		}

		return nil, ErrInternal
	}

	return &userpb.RegisterResponse{
		Status:    StatusSuccess,
		Timestamp: time.Now().Unix(),
		Result: &userpb.RegistrationResult{
			Username: userDetails.Name,
			UserID:   userDetails.ID,
			Password: req.Password,
		},
	}, nil
}
