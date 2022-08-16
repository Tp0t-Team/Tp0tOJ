package utils

import (
	"time"
)

type RankCache interface {
	SetCalculator(calculator ScoreCalculator)
	GetRank() []RankItem
	GetUserScore(userId uint64) uint64
	AddUser(userId uint64)
	AddChallenge(challengeId uint64, originScore uint64)
	Submit(userId uint64, challengeId uint64, stamp time.Time) error
	GetCurrentScores() map[uint64]uint64
	WarmUp() error
}

type RankItem struct {
	UserId uint64
	Score  uint64
}

var Cache RankCache
