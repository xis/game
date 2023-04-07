package domain

import "context"

type Leaderboard struct {
	UserScores []UserScore
}

type UserScore struct {
	UserID   string
	Username string
	Score    float64
}

//go:generate mockery --name UserScoreRepository --structname MockUserScoreRepository --outpkg mocks --filename user_score_repository_mock.go --output ./mocks/. --with-expecter
type UserScoreRepository interface {
	GetUserTopScore(ctx context.Context, userID string) (UserScore, error)
	UpdateUserTopScore(ctx context.Context, userID string, score float64) error
	GetLeaderboard(ctx context.Context) (Leaderboard, error)
}
