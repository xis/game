package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"game/internal/domain"
	"game/internal/domain/mocks"
)

type LeaderboardServiceTestSuite struct {
	suite.Suite

	service *leaderboardService

	mockUserRepository      *mocks.MockUserRepository
	mockUserScoreRepository *mocks.MockUserScoreRepository
}

func TestLeaderboardServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LeaderboardServiceTestSuite))
}

func (suite *LeaderboardServiceTestSuite) SetupTest() {
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.T())
	suite.mockUserScoreRepository = mocks.NewMockUserScoreRepository(suite.T())

	suite.service = NewLeaderboardService(LeaderboardServiceDependencies{
		UserRepository:      suite.mockUserRepository,
		UserScoreRepository: suite.mockUserScoreRepository,
	})
}

func (suite *LeaderboardServiceTestSuite) TestGetLeaderboard() {
	suite.mockUserScoreRepository.
		EXPECT().
		GetLeaderboard(mock.Anything).
		Return(domain.Leaderboard{}, nil)

	_, err := suite.service.GetLeaderboard(context.Background())
	suite.NoError(err)
}

func (suite *LeaderboardServiceTestSuite) TestGetLeaderboard_RepositoryFailed() {
	suite.mockUserScoreRepository.
		EXPECT().
		GetLeaderboard(mock.Anything).
		Return(domain.Leaderboard{}, domain.ErrInternal)

	_, err := suite.service.GetLeaderboard(context.Background())
	suite.Error(err)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(true, nil)

	suite.mockUserScoreRepository.
		EXPECT().
		GetUserTopScore(mock.Anything, "user-id").
		Return(domain.UserScore{}, domain.ErrResourceNotFound)

	suite.mockUserScoreRepository.
		EXPECT().
		UpdateUserTopScore(mock.Anything, "user-id", float64(10)).
		Return(nil)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.NoError(err)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore_GetUserTopScoreFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(true, nil)

	suite.mockUserScoreRepository.
		EXPECT().
		GetUserTopScore(mock.Anything, "user-id").
		Return(domain.UserScore{}, domain.ErrInternal)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.Error(err)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore_UpdateUserTopScoreFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(true, nil)

	suite.mockUserScoreRepository.
		EXPECT().
		GetUserTopScore(mock.Anything, "user-id").
		Return(domain.UserScore{}, domain.ErrResourceNotFound)

	suite.mockUserScoreRepository.
		EXPECT().
		UpdateUserTopScore(mock.Anything, "user-id", float64(10)).
		Return(domain.ErrInternal)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.Error(err)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore_UpdateUserTopScoreSkipped() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(true, nil)

	suite.mockUserScoreRepository.
		EXPECT().
		GetUserTopScore(mock.Anything, "user-id").
		Return(domain.UserScore{
			Score: 20,
		}, nil)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.NoError(err)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore_UserNotFound() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(false, nil)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.ErrorIs(err, domain.ErrResourceNotFound)
}

func (suite *LeaderboardServiceTestSuite) TestSubmitUserScore_CheckExistsByIDFailed() {
	suite.mockUserRepository.
		EXPECT().
		CheckExistsByID(mock.Anything, "user-id").
		Return(false, domain.ErrInternal)

	err := suite.service.SubmitUserScore(context.Background(), "user-id", 10)
	suite.Error(err)
}
