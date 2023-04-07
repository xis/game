package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"game/internal/domain/mocks"
)

type UnaryInterceptorTestSuite struct {
	suite.Suite

	interceptor *UnaryInterceptor

	mockTokenManager *mocks.MockTokenManager
}

func TestUnaryInterceptorTestSuite(t *testing.T) {
	suite.Run(t, new(UnaryInterceptorTestSuite))
}

func (suite *UnaryInterceptorTestSuite) SetupTest() {
	suite.mockTokenManager = mocks.NewMockTokenManager(suite.T())

	suite.interceptor = NewUnaryInterceptor(UnaryInterceptorDependencies{
		TokenManager: suite.mockTokenManager,
		AuthorizedMethodNames: []string{
			"some-method",
			"a-different-method",
		},
	})
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor() {
	suite.mockTokenManager.
		EXPECT().
		ExtractUserID(mock.Anything, "token").
		Return("user-id", nil)

	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	expectedUserID := "user-id"

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		userID := ctx.Value(ContextKeyUserID).(string)

		suite.Equal(expectedUserID, userID)

		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"Authorization": "Bearer token",
	}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.NoError(err)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_NotAuthorizedMethod() {
	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "not-authorizd",
	}, unaryHandler)

	suite.NoError(err)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_InvalidMetadata() {
	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}

	_, err := suite.interceptor.Intercept(context.Background(), exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.ErrorIs(err, ErrInvalidMetadata)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_NoAuthorizationHeader() {
	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.ErrorIs(err, ErrUnauthenticated)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_InvalidAuthorizationHeader() {
	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "invalid",
	}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.ErrorIs(err, ErrUnauthenticated)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_EmptyToken() {
	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "Bearer ",
	}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.ErrorIs(err, ErrUnauthenticated)
}

func (suite *UnaryInterceptorTestSuite) TestUnaryInterceptor_InvalidToken() {
	someError := errors.New("some-error")

	suite.mockTokenManager.
		EXPECT().
		ExtractUserID(mock.Anything, "token").
		Return("", someError)

	exampleReq := struct {
		Score float64
	}{
		Score: 86,
	}

	expectedUserID := "user-id"

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		userID := ctx.Value(ContextKeyUserID).(string)

		suite.Equal(expectedUserID, userID)

		return req, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"Authorization": "Bearer token",
	}))

	_, err := suite.interceptor.Intercept(ctx, exampleReq, &grpc.UnaryServerInfo{
		FullMethod: "some-method",
	}, unaryHandler)

	suite.ErrorIs(err, someError)
}
