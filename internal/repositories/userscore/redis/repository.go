package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"game/internal/domain"
)

const (
	leaderboardKey = "leaderboard"
)

type RedisUserScoreRepositoryDependencies struct {
	Client *redis.Client

	UserRepository domain.UserRepository
}

type RedisUserScoreRepository struct {
	client         *redis.Client
	userRepository domain.UserRepository
}

func NewRedisUserScoreRepository(deps RedisUserScoreRepositoryDependencies) *RedisUserScoreRepository {
	return &RedisUserScoreRepository{
		client:         deps.Client,
		userRepository: deps.UserRepository,
	}
}

func (repo *RedisUserScoreRepository) GetUserTopScore(ctx context.Context, userID string) (domain.UserScore, error) {
	score, err := repo.client.ZScore(ctx, leaderboardKey, userID).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.UserScore{}, domain.ErrResourceNotFound
		}

		return domain.UserScore{}, err
	}

	return domain.UserScore{
		UserID: userID,
		Score:  score,
	}, nil
}

func (repo *RedisUserScoreRepository) UpdateUserTopScore(ctx context.Context, userID string, score float64) error {
	_, err := repo.client.ZAdd(ctx, leaderboardKey, &redis.Z{
		Score:  score,
		Member: userID,
	}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (repo *RedisUserScoreRepository) GetLeaderboard(ctx context.Context) (domain.Leaderboard, error) {
	userScores, err := repo.client.ZRevRangeWithScores(ctx, leaderboardKey, 0, -1).Result()
	if err != nil {
		return domain.Leaderboard{}, err
	}

	leaderboard := domain.Leaderboard{
		UserScores: make([]domain.UserScore, len(userScores)),
	}

	var userIDs []string

	for _, userScore := range userScores {
		userID, ok := userScore.Member.(string)
		if !ok {
			return domain.Leaderboard{}, fmt.Errorf("%w, invalid user id type: %T", domain.ErrInternal, userScore.Member)
		}

		userIDs = append(userIDs, userID)
	}

	users, err := repo.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return domain.Leaderboard{}, err
	}

	userByID := make(map[string]domain.User)

	for _, user := range users {
		userByID[user.ID] = user
	}

	for i, userScore := range userScores {
		user, ok := userByID[userScore.Member.(string)]
		if !ok {
			return domain.Leaderboard{}, fmt.Errorf("%w, user not found: %s", domain.ErrInternal, userScore.Member.(string))
		}

		leaderboard.UserScores[i] = domain.UserScore{
			UserID:   user.ID,
			Username: user.Name,
			Score:    userScore.Score,
		}
	}

	return leaderboard, nil
}
