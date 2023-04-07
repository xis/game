package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"game/internal/domain"
	leaderboardpb "game/internal/proto/leaderboard/proto"
	"game/internal/services"
)

var (
	ErrInvalidUserID = status.New(codes.InvalidArgument, "invalid user id").Err()
	ErrInvalidScore  = status.New(codes.InvalidArgument, "invalid score").Err()
)

type LeaderboardControllerDependencies struct {
	LeaderboardService services.LeaderboardService

	Logger *logrus.Logger
}

type leaderboardController struct {
	leaderboardpb.UnimplementedLeaderboardServiceServer

	leaderboardService services.LeaderboardService

	logger *logrus.Logger
}

func NewLeaderboardController(deps LeaderboardControllerDependencies) *leaderboardController {
	return &leaderboardController{
		leaderboardService: deps.LeaderboardService,
		logger:             deps.Logger,
	}
}

func (controller *leaderboardController) GetLeaderboard(ctx context.Context, request *leaderboardpb.GetLeaderboardRequest) (*leaderboardpb.GetLeaderboardResponse, error) {
	controller.logger.Info("get leaderboard request has been received")

	leaderboard, err := controller.leaderboardService.GetLeaderboard(ctx)
	if err != nil {
		controller.logger.
			WithError(err).
			Error("failed to get leaderboard")

		return nil, ErrInternal
	}

	var results []*leaderboardpb.UserScore

	for _, userScore := range leaderboard.UserScores {
		results = append(results, &leaderboardpb.UserScore{
			Username: userScore.Username,
			Score:    userScore.Score,
			UserID:   userScore.UserID,
		})
	}

	return &leaderboardpb.GetLeaderboardResponse{
		Status:    StatusSuccess,
		Timestamp: time.Now().Unix(),
		Results:   results,
	}, nil
}

func (controller *leaderboardController) SubmitUserScore(ctx context.Context, request *leaderboardpb.SubmitUserScoreRequest) (*leaderboardpb.SubmitUserScoreResponse, error) {
	controller.logger.Info("submit user score request has been received")

	if request.Score <= 0 {
		return nil, ErrInvalidScore
	}

	userID, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok {
		return nil, ErrInvalidUserID
	}

	err := controller.leaderboardService.SubmitUserScore(ctx, userID, request.Score)
	if err != nil {
		controller.logger.
			WithError(err).
			WithField("user_id", userID).
			Error("failed to submit user score")

		if errors.Is(err, domain.ErrResourceNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternal
	}

	return &leaderboardpb.SubmitUserScoreResponse{
		Status:    StatusSuccess,
		Timestamp: time.Now().Unix(),
	}, nil
}
