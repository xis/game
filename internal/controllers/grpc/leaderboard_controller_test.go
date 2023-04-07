package grpc

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"game/internal/domain"
	leaderboardpb "game/internal/proto/leaderboard/proto"
	"game/internal/services/mocks"
)

type LeaderboardControllerTestSuite struct {
	suite.Suite

	controller *leaderboardController

	mockLeaderboardService *mocks.MockLeaderboardService
}

func TestLeaderboardControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LeaderboardControllerTestSuite))
}

func (suite *LeaderboardControllerTestSuite) SetupTest() {
	suite.mockLeaderboardService = mocks.NewMockLeaderboardService(suite.T())

	suite.controller = NewLeaderboardController(LeaderboardControllerDependencies{
		LeaderboardService: suite.mockLeaderboardService,

		Logger: logrus.New(),
	})
}

func (suite *LeaderboardControllerTestSuite) TestGetLeaderboard() {
	suite.mockLeaderboardService.
		EXPECT().
		GetLeaderboard(mock.Anything).
		Return(domain.Leaderboard{
			UserScores: []domain.UserScore{
				{
					UserID:   "user-id",
					Username: "username",
					Score:    86,
				},
				{
					UserID:   "user-id-2",
					Username: "username-2",
					Score:    82,
				},
			},
		}, nil)

	result, err := suite.controller.GetLeaderboard(context.Background(), &leaderboardpb.GetLeaderboardRequest{})
	suite.NoError(err)

	expectedResult := &leaderboardpb.GetLeaderboardResponse{
		Status: StatusSuccess,
		Results: []*leaderboardpb.UserScore{
			{
				UserID:   "user-id",
				Username: "username",
				Score:    86,
			},
			{
				UserID:   "user-id-2",
				Username: "username-2",
				Score:    82,
			},
		},
	}

	suite.Equal(expectedResult.Status, result.Status)
	suite.Equal(expectedResult.Results, result.Results)
	suite.NotEmpty(result.Timestamp)
}

func (suite *LeaderboardControllerTestSuite) TestGetLeaderboard_ServiceFailed() {
	suite.mockLeaderboardService.
		EXPECT().
		GetLeaderboard(mock.Anything).
		Return(domain.Leaderboard{}, domain.ErrInternal)

	result, err := suite.controller.GetLeaderboard(context.Background(), &leaderboardpb.GetLeaderboardRequest{})
	suite.ErrorIs(err, ErrInternal)
	suite.Empty(result)
}

func (suite *LeaderboardControllerTestSuite) TestSubmitUserScore() {
	suite.mockLeaderboardService.
		EXPECT().
		SubmitUserScore(mock.Anything, "user-id", float64(86)).
		Return(nil)

	ctx := context.WithValue(context.Background(), ContextKeyUserID, "user-id")

	result, err := suite.controller.SubmitUserScore(ctx, &leaderboardpb.SubmitUserScoreRequest{
		Score: 86,
	})
	suite.NoError(err)

	expectedResult := &leaderboardpb.SubmitUserScoreResponse{
		Status: StatusSuccess,
	}

	suite.Equal(expectedResult.Status, result.Status)
	suite.NotEmpty(result.Timestamp)
}

func (suite *LeaderboardControllerTestSuite) TestSubmitUserScore_ServiceFailed() {
	suite.mockLeaderboardService.
		EXPECT().
		SubmitUserScore(mock.Anything, "user-id", float64(86)).
		Return(domain.ErrInternal)

	ctx := context.WithValue(context.Background(), ContextKeyUserID, "user-id")

	result, err := suite.controller.SubmitUserScore(ctx, &leaderboardpb.SubmitUserScoreRequest{
		Score: 86,
	})
	suite.ErrorIs(err, ErrInternal)
	suite.Empty(result)
}

func (suite *LeaderboardControllerTestSuite) TestSubmitUserScore_NoUserID() {
	result, err := suite.controller.SubmitUserScore(context.Background(), &leaderboardpb.SubmitUserScoreRequest{
		Score: 86,
	})
	suite.ErrorIs(err, ErrInvalidUserID)
	suite.Empty(result)
}

func (suite *LeaderboardControllerTestSuite) TestSubmitUserScore_InvalidScore() {
	result, err := suite.controller.SubmitUserScore(context.Background(), &leaderboardpb.SubmitUserScoreRequest{
		Score: -1,
	})
	suite.ErrorIs(err, ErrInvalidScore)
	suite.Empty(result)
}

func (suite *LeaderboardControllerTestSuite) TestSubmitUserScore_ResourceNotFound() {
	suite.mockLeaderboardService.
		EXPECT().
		SubmitUserScore(mock.Anything, "user-id", float64(86)).
		Return(domain.ErrResourceNotFound)

	ctx := context.WithValue(context.Background(), ContextKeyUserID, "user-id")

	result, err := suite.controller.SubmitUserScore(ctx, &leaderboardpb.SubmitUserScoreRequest{
		Score: 86,
	})
	suite.ErrorIs(err, domain.ErrResourceNotFound)
	suite.Empty(result)
}
