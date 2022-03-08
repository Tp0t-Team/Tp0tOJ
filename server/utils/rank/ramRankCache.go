package rank

import (
	"encoding/json"
	"errors"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils"
	"sort"
	"strconv"
	"sync"
	"time"
)

type RAMRankCache struct {
	calculator     utils.ScoreCalculator
	mutex          sync.RWMutex
	rank           []uint64
	challengeScore map[uint64]uint64
	challengeSolve map[uint64][]uint64
	userScore      map[uint64]uint64
	userTime       map[uint64]time.Time
}

func init() {
	//calculator need to be init before usage
	utils.Cache = &RAMRankCache{rank: []uint64{}, challengeScore: map[uint64]uint64{}, challengeSolve: map[uint64][]uint64{}, userScore: map[uint64]uint64{}, userTime: map[uint64]time.Time{}}
}

func (cache *RAMRankCache) SetCalculator(calculator utils.ScoreCalculator) {
	cache.calculator = calculator
}

func (cache *RAMRankCache) GetRank() []utils.RankItem {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	ret := []utils.RankItem{}
	for _, user := range cache.rank {
		ret = append(ret, utils.RankItem{
			UserId: user,
			Score:  cache.userScore[user],
		})
	}
	return ret
}

func (cache *RAMRankCache) GetUserScore(userId uint64) uint64 {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	score, ok := cache.userScore[userId]
	if !ok {
		return 0
	}
	return score
}

func (cache *RAMRankCache) AddUser(userId uint64) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.userScore[userId] = 0
	cache.userTime[userId] = time.Time{}
}

func (cache *RAMRankCache) AddChallenge(challengeId uint64, originScore uint64) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.challengeScore[challengeId] = originScore
	cache.challengeSolve[challengeId] = []uint64{}
}

func (cache *RAMRankCache) Submit(userId uint64, challengeId uint64, stamp time.Time) error {
	if err := cache.submitImpl(userId, challengeId, stamp); err != nil {
		return err
	}
	cache.refreshRank()
	return nil
}
func (cache *RAMRankCache) submitImpl(userId uint64, challengeId uint64, stamp time.Time) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	oldScore, exist := cache.challengeScore[challengeId]
	if !exist {
		return errors.New("unexist challenge")
	}
	count := uint64(len(cache.challengeSolve[challengeId])) + 1
	newScore := cache.calculator.GetScore(challengeId, count)

	cache.userTime[userId] = stamp
	cache.challengeSolve[challengeId] = append(cache.challengeSolve[challengeId], userId)
	cache.userScore[userId] += cache.calculator.GetIncrementScore(oldScore, count-1)
	cache.challengeScore[challengeId] = newScore

	for index, user := range cache.challengeSolve[challengeId] {
		cache.userScore[user] -= cache.calculator.GetDeltaScoreForUser(oldScore, newScore, uint64(index))
	}
	return nil
}

type ScoreItem struct {
	userId uint64
	score  uint64
	stamp  time.Time
}

type ScoreItems = []*ScoreItem

type RankSorter struct {
	ScoreItems
}

func (items RankSorter) Len() int {
	return len(items.ScoreItems)
}

func (items RankSorter) Swap(i, j int) {
	items.ScoreItems[i], items.ScoreItems[j] = items.ScoreItems[j], items.ScoreItems[i]
}

func (items RankSorter) Less(i, j int) bool {
	return items.ScoreItems[i].score > items.ScoreItems[j].score || (items.ScoreItems[i].score == items.ScoreItems[j].score && items.ScoreItems[i].stamp.Before(items.ScoreItems[j].stamp))
}

func (cache *RAMRankCache) refreshRank() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	rank := ScoreItems{}
	for user, score := range cache.userScore {
		rank = append(rank, &ScoreItem{
			userId: user,
			score:  score,
			stamp:  cache.userTime[user],
		})
	}
	sort.Sort(RankSorter{rank})
	var rankSaved []uint64
	for _, item := range rank {
		rankSaved = append(rankSaved, item.userId)
	}
	cache.rank = rankSaved
}

func (cache *RAMRankCache) WarmUp() error {
	challenges := resolvers.FindAllChallenges()
	if challenges == nil {
		return errors.New("challenges equals nil")
	}
	for _, challenge := range challenges {
		if challenge.State == "disabled" {
			continue
		}
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			return err
		}
		var score uint64
		score, err = strconv.ParseUint(config.Score.BaseScore, 10, 64)
		if err != nil {
			return err
		}
		cache.AddChallenge(challenge.ChallengeId, score)
	}
	users := resolvers.FindAllUser()
	//if users == nil {
	//	return errors.New("users equals nil")
	//}
	if users != nil {
		for _, user := range users {
			cache.AddUser(user.UserId)
		}
	}
	submits := resolvers.FindSubmitCorrectSorted()
	//if submits == nil {
	//	return errors.New("users equals nil")
	//}
	if submits != nil {
		for _, submit := range submits {
			err := cache.submitImpl(submit.UserId, submit.ChallengeId, submit.SubmitTime)
			if err != nil {
				return err
			}
		}
	}
	cache.refreshRank()
	return nil
}
