package utils

type ScoreCalculator interface {
	//GetScore(challengeId uint64, count uint64, dynamic bool) uint64
	GetScore(baseScore uint64, count uint64, dynamic bool) uint64
	GetIncrementScore(score uint64, index uint64) uint64
	GetDeltaScoreForUser(oldScore uint64, newScore uint64, index uint64) uint64
}
