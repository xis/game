package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"game/internal/domain"
)

var (
	ErrUnauthenticated = status.New(codes.Unauthenticated, "unauthenticated").Err()
	ErrInvalidMetadata = status.New(codes.InvalidArgument, "invalid metadata").Err()
)

type ContextKey string

const (
	ContextKeyUserID ContextKey = "user_id"
)

type UnaryInterceptorDependencies struct {
	TokenManager          domain.TokenManager
	AuthorizedMethodNames []string
}

type UnaryInterceptor struct {
	tokenManager          domain.TokenManager
	authorizedMethodNames map[string]struct{}
}

func NewUnaryInterceptor(
	deps UnaryInterceptorDependencies,
) *UnaryInterceptor {
	authorizedMethodNames := make(map[string]struct{})

	for _, methodName := range deps.AuthorizedMethodNames {
		authorizedMethodNames[methodName] = struct{}{}
	}

	return &UnaryInterceptor{
		tokenManager:          deps.TokenManager,
		authorizedMethodNames: authorizedMethodNames,
	}
}

func (interceptor *UnaryInterceptor) Intercept(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {

	if interceptor.isMethodAuthorized(info.FullMethod) {
		userID, err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, ContextKeyUserID, userID)
	}

	return handler(ctx, req)
}

func (interceptor *UnaryInterceptor) isMethodAuthorized(method string) bool {
	if _, ok := interceptor.authorizedMethodNames[method]; ok {
		return true
	}

	return false
}

func (interceptor *UnaryInterceptor) authorize(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrInvalidMetadata
	}

	authorizationHeaderValue := md["authorization"]
	if len(authorizationHeaderValue) == 0 {
		return "", ErrUnauthenticated
	}

	headerParts := strings.Split(authorizationHeaderValue[0], " ")
	if len(headerParts) != 2 {
		return "", ErrUnauthenticated
	}

	token := headerParts[1]

	if token == "" {
		return "", ErrUnauthenticated
	}

	userID, err := interceptor.tokenManager.ExtractUserID(ctx, token)
	if err != nil {
		return "", err
	}

	return userID, nil
}
