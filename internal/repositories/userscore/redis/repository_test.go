package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"game/internal/domain"
	"game/internal/domain/mocks"
)

type RedisUserScoreRepositoryTestSuite struct {
	suite.Suite

	repository *RedisUserScoreRepository

	redisMock          redismock.ClientMock
	mockUserRepository *mocks.MockUserRepository
}

func TestRedisUserScoreRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RedisUserScoreRepositoryTestSuite))
}

func (suite *RedisUserScoreRepositoryTestSuite) SetupTest() {
	db, mock := redismock.NewClientMock()

	suite.redisMock = mock
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.T())

	suite.repository = NewRedisUserScoreRepository(RedisUserScoreRepositoryDependencies{
		Client:         db,
		UserRepository: suite.mockUserRepository,
	})
}

func (suite *RedisUserScoreRepositoryTestSuite) TearDownTest() {
	suite.NoError(suite.redisMock.ExpectationsWereMet())
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetUserTopScore() {
	suite.redisMock.
		ExpectZScore("leaderboard", "user-id").
		SetVal(100)

	_, err := suite.repository.GetUserTopScore(context.Background(), "user-id")
	suite.NoError(err)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetUserTopScore_ZScoreFailed() {
	someError := errors.New("some error")

	suite.redisMock.
		ExpectZScore("leaderboard", "user-id").
		SetErr(someError)

	_, err := suite.repository.GetUserTopScore(context.Background(), "user-id")
	suite.ErrorIs(err, someError)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetUserTopScore_ZScoreNotFound() {
	suite.redisMock.
		ExpectZScore("leaderboard", "user-id").
		RedisNil()

	_, err := suite.repository.GetUserTopScore(context.Background(), "user-id")
	suite.ErrorIs(err, domain.ErrResourceNotFound)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestUpdateUserTopScore() {
	suite.redisMock.
		ExpectZAdd("leaderboard", &redis.Z{
			Score:  900,
			Member: "user-id",
		}).
		SetVal(900)

	err := suite.repository.UpdateUserTopScore(context.Background(), "user-id", 900)
	suite.NoError(err)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestUpdateUserTopScore_ZAddFailed() {
	someError := errors.New("some error")

	suite.redisMock.
		ExpectZAdd("leaderboard", &redis.Z{
			Score:  900,
			Member: "user-id",
		}).
		SetErr(someError)

	err := suite.repository.UpdateUserTopScore(context.Background(), "user-id", 900)
	suite.ErrorIs(err, someError)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetLeaderboard() {
	suite.redisMock.
		ExpectZRevRangeWithScores("leaderboard", 0, -1).
		SetVal([]redis.Z{
			{
				Score:  900,
				Member: "user-id-1",
			},
			{
				Score:  800,
				Member: "user-id-2",
			},
		})

	suite.mockUserRepository.
		EXPECT().
		GetUsersByIDs(mock.Anything, []string{"user-id-1", "user-id-2"}).
		Return([]domain.User{
			{
				ID:           "user-id-1",
				Name:         "user-1",
				PasswordHash: "password-hash",
			},
			{
				ID:           "user-id-2",
				Name:         "user-2",
				PasswordHash: "password-hash",
			},
		}, nil)

	_, err := suite.repository.GetLeaderboard(context.Background())
	suite.NoError(err)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetLeaderboard_ZRevRangeWithScoresFailed() {
	someError := errors.New("some error")

	suite.redisMock.
		ExpectZRevRangeWithScores("leaderboard", 0, -1).
		SetErr(someError)

	_, err := suite.repository.GetLeaderboard(context.Background())
	suite.ErrorIs(err, someError)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetLeaderboard_GetUsersByIDsFailed() {
	someError := errors.New("some error")

	suite.redisMock.
		ExpectZRevRangeWithScores("leaderboard", 0, -1).
		SetVal([]redis.Z{
			{
				Score:  900,
				Member: "user-id-1",
			},
			{
				Score:  800,
				Member: "user-id-2",
			},
		})

	suite.mockUserRepository.
		EXPECT().
		GetUsersByIDs(mock.Anything, []string{"user-id-1", "user-id-2"}).
		Return(nil, someError)

	_, err := suite.repository.GetLeaderboard(context.Background())
	suite.ErrorIs(err, someError)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetLeaderboard_GetUsersByIDsNotFound() {
	suite.redisMock.
		ExpectZRevRangeWithScores("leaderboard", 0, -1).
		SetVal([]redis.Z{
			{
				Score:  900,
				Member: "user-id-1",
			},
			{
				Score:  800,
				Member: "user-id-2",
			},
		})

	suite.mockUserRepository.
		EXPECT().
		GetUsersByIDs(mock.Anything, []string{"user-id-1", "user-id-2"}).
		Return(nil, nil)

	_, err := suite.repository.GetLeaderboard(context.Background())
	suite.ErrorIs(err, domain.ErrInternal)
}

func (suite *RedisUserScoreRepositoryTestSuite) TestGetLeaderboard_GetUsersByIDsEmpty() {
	suite.redisMock.
		ExpectZRevRangeWithScores("leaderboard", 0, -1).
		SetVal([]redis.Z{
			{
				Score:  900,
				Member: "user-id-1",
			},
			{
				Score:  800,
				Member: "user-id-2",
			},
		})

	suite.mockUserRepository.
		EXPECT().
		GetUsersByIDs(mock.Anything, []string{"user-id-1", "user-id-2"}).
		Return([]domain.User{}, nil)

	_, err := suite.repository.GetLeaderboard(context.Background())
	suite.ErrorIs(err, domain.ErrInternal)
}
