package jwt

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type JWTTokenManagerTestSuite struct {
	suite.Suite

	tokenManager *JWTTokenManager
}

func TestJWTTokenManagerTestSuite(t *testing.T) {
	suite.Run(t, new(JWTTokenManagerTestSuite))
}

func (suite *JWTTokenManagerTestSuite) SetupTest() {
	suite.tokenManager = NewJWTTokenManager(JWTTokenCreatorDependencies{
		SecretKey: "secret-key",
	})
}

func (suite *JWTTokenManagerTestSuite) TestCreate() {
	token, err := suite.tokenManager.Create(context.Background(), "user-id")
	suite.NoError(err)
	suite.NotEmpty(token)
}

func (suite *JWTTokenManagerTestSuite) TestExtractUserID() {
	token, err := suite.tokenManager.Create(context.Background(), "user-id")
	suite.NoError(err)
	suite.NotEmpty(token)

	userID, err := suite.tokenManager.ExtractUserID(context.Background(), token)
	suite.NoError(err)
	suite.Equal("user-id", userID)
}

func (suite *JWTTokenManagerTestSuite) TestExtractUserID_InvalidToken() {
	_, err := suite.tokenManager.ExtractUserID(context.Background(), "invalid-token")
	suite.Error(err)
}

func (suite *JWTTokenManagerTestSuite) TestExtractUserID_ExpiredToken() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: "user-id",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 8000,
		},
	})

	tokenString, err := token.SignedString([]byte(suite.tokenManager.secretKey))
	suite.NoError(err)

	_, err = suite.tokenManager.ExtractUserID(context.Background(), tokenString)
	suite.ErrorContains(err, "expired")
}

func (suite *JWTTokenManagerTestSuite) TestExtractUserID_InvalidSignature() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: "user-id",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 8000,
		},
	})

	tokenString, err := token.SignedString([]byte("invalid-secret-key"))
	suite.NoError(err)

	_, err = suite.tokenManager.ExtractUserID(context.Background(), tokenString)
	suite.ErrorContains(err, "signature is invalid")
}

func (suite *JWTTokenManagerTestSuite) TestExtractUserID_NoUserID() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 8000000000,
		},
	})

	tokenString, err := token.SignedString([]byte(suite.tokenManager.secretKey))
	suite.NoError(err)

	_, err = suite.tokenManager.ExtractUserID(context.Background(), tokenString)
	suite.ErrorContains(err, "invalid token")
}
