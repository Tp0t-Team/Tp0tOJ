package utils

import (
	"encoding/json"
	"math"
	"server/services/database/resolvers"
	"server/services/types"
	"strconv"
)

type ScoreCalculator interface {
	GetScore(challengeId uint64, count uint64) uint64
	GetIncrementScore(score uint64, index uint64) uint64
	GetDeltaScoreForUser(oldScore uint64, newScore uint64, index uint64) uint64
}

type BasicScoreCalculator struct{}

func (c *BasicScoreCalculator) curve(baseScore uint64, count uint64) uint64 {
	if count == 0 {
		return baseScore
	}
	count -= 1
	coefficient := 1.8414 / float64(Configure.Challenge.HalfLife) * float64(count)
	result := math.Floor(float64(baseScore) / (coefficient + math.Exp(-coefficient)))
	return uint64(result)
}

func (c *BasicScoreCalculator) GetScore(challengeId uint64, count uint64) uint64 {
	// step 1: from challengeId, get curve parameters
	challenge, err := resolvers.FindChallengeById(challengeId)
	if err != nil {
		return 0
	}
	var config types.ChallengeConfig
	err = json.Unmarshal([]byte(challenge.Configuration), &config)
	if err != nil {
		return 0
	}
	// step 2: use curve(), parameters & count, calc the score
	score, err := strconv.ParseUint(config.Score.BaseScore, 10, 64)
	if err != nil {
		return 0
	}
	if config.Score.Dynamic {
		return c.curve(score, count)
	} else {
		return score
	}
}

func (c *BasicScoreCalculator) GetIncrementScore(score uint64, index uint64) uint64 {
	if index == 0 {
		return uint64(math.Floor(float64(score) * (1 + Configure.Challenge.FirstBloodReward)))
	} else if index == 1 {
		return uint64(math.Floor(float64(score) * (1 + Configure.Challenge.SecondBloodReward)))
	} else if index == 2 {
		return uint64(math.Floor(float64(score) * (1 + Configure.Challenge.ThirdBloodReward)))
	} else {
		return score
	}
}

func (c *BasicScoreCalculator) GetDeltaScoreForUser(oldScore uint64, newScore uint64, index uint64) uint64 {
	if index == 0 {
		return uint64(math.Floor(float64(oldScore)*(1+Configure.Challenge.FirstBloodReward))) - uint64(math.Floor(float64(newScore)*(1+Configure.Challenge.FirstBloodReward)))
	} else if index == 1 {
		return uint64(math.Floor(float64(oldScore)*(1+Configure.Challenge.SecondBloodReward))) - uint64(math.Floor(float64(newScore)*(1+Configure.Challenge.SecondBloodReward)))
	} else if index == 2 {
		return uint64(math.Floor(float64(oldScore)*(1+Configure.Challenge.ThirdBloodReward))) - uint64(math.Floor(float64(newScore)*(1+Configure.Challenge.ThirdBloodReward)))
	} else {
		return oldScore - newScore
	}
}
