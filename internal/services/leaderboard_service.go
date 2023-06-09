package services

import (
	"context"

	"game/internal/domain"
)

//go:generate mockery --name LeaderboardService --structname MockLeaderboardService --outpkg mocks --filename leaderboard_service_mock.go --output ./mocks/. --with-expecter
type LeaderboardService interface {
	GetLeaderboard(ctx context.Context) (domain.Leaderboard, error)
	SubmitUserScore(ctx context.Context, userID string, score float64) error
}

type LeaderboardServiceDependencies struct {
	UserScoreRepository domain.UserScoreRepository
	UserRepository      domain.UserRepository
}

type leaderboardService struct {
	userScoreRepository domain.UserScoreRepository
	userRepository      domain.UserRepository
}

func NewLeaderboardService(deps LeaderboardServiceDependencies) *leaderboardService {
	return &leaderboardService{
		userScoreRepository: deps.UserScoreRepository,
		userRepository:      deps.UserRepository,
	}
}

func (service *leaderboardService) GetLeaderboard(ctx context.Context) (domain.Leaderboard, error) {
	leaderboard, err := service.userScoreRepository.GetLeaderboard(ctx)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	return leaderboard, nil
}

func (service *leaderboardService) SubmitUserScore(ctx context.Context, userID string, score float64) error {
	exists, err := service.userRepository.CheckExistsByID(ctx, userID)
	if err != nil {
		return err
	}

	if !exists {
		return domain.ErrResourceNotFound
	}

	userTopScore, err := service.userScoreRepository.GetUserTopScore(ctx, userID)
	if err != nil && err != domain.ErrResourceNotFound {
		return err
	}

	if userTopScore.Score > score {
		return nil
	}

	err = service.userScoreRepository.UpdateUserTopScore(ctx, userID, score)
	if err != nil {
		return err
	}

	return nil
}
