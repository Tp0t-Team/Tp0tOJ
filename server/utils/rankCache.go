package utils

import (
	"time"
)

type ChartCurve struct {
	Id    uint64   `json:"id"`
	Name  string   `json:"name"`
	Score []uint64 `json:"score"`
}

type ChartData struct {
	X    []int64       `json:"x"`
	Y    []*ChartCurve `json:"y"`
	TopN uint64        `json:"topN"`
}

type RankCache interface {
	SetCalculator(calculator ScoreCalculator)
	GetRank() []RankItem
	GetUserScore(userId uint64) uint64
	//AddUser(userId uint64)
	MutateUser(userId uint64, state bool)
	//AddChallenge(challengeId uint64, originScore uint64)
	MutateChallenge(challengeId uint64, state bool, dynamic bool, baseScore uint64)
	Submit(userId uint64, challengeId uint64, stamp time.Time) error
	GetCurrentScores() map[uint64]uint64
	//WarmUp() error
	Load(filename string) error

	Chart(topN uint64) *ChartData
}

type RankItem struct {
	UserId uint64
	Score  uint64
}

var Cache RankCache
