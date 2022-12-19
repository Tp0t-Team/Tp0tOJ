package rank

import (
	"encoding/json"
	"errors"
	"github.com/ugorji/go/codec"
	"io"
	"os"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils"
	"strconv"
	"sync"
	"time"
)

type ChallengeFrame struct {
	Solved uint64
}

type UserFrame struct {
	Solved []uint64
	Blood  map[uint64]int
}

type TimeFrame struct {
	Challenges map[uint64]*ChallengeFrame
	Users      map[uint64]*UserFrame
	Stamp      time.Time
}

type TimelineRankCache struct {
	calculator       utils.ScoreCalculator
	mutex            sync.RWMutex
	rankNodes        []ScoreItems
	frames           []TimeFrame
	challengeState   map[uint64]bool
	challengeDynamic map[uint64]bool
	challengeScore   map[uint64]uint64
	userState        map[uint64]bool

	fileCache *FileCache
}

type FileCache struct {
	file     *os.File
	filename string
	encoder  *codec.Encoder
	decoder  *codec.Decoder
}

func NewFileCache(filename string) (*FileCache, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	handle := &codec.CborHandle{}
	ret := &FileCache{
		file:     file,
		filename: filename,
		encoder:  codec.NewEncoder(file, handle),
		decoder:  codec.NewDecoder(file, handle),
	}
	return ret, nil
}

func (cache *FileCache) Load() ([]TimeFrame, error) {
	ret := []TimeFrame{}
	for {
		item := TimeFrame{}
		err := cache.decoder.Decode(&item)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		ret = append(ret, item)
	}
	return ret, nil
}

func (cache *FileCache) Append(frame *TimeFrame) error {
	err := cache.encoder.Encode(*frame)
	if err != nil {
		return err
	}
	return cache.file.Sync()
}

func init() {
	//calculator need to be init before usage
	// TODO:
	//utils.Cache = &RAMRankCache{rank: []uint64{}, challengeScore: map[uint64]uint64{}, challengeSolve: map[uint64][]uint64{}, userScore: map[uint64]uint64{}, userTime: map[uint64]time.Time{}}
}

func (cache *TimelineRankCache) SetCalculator(calculator utils.ScoreCalculator) {
	cache.calculator = calculator
}

func (cache *TimelineRankCache) GetRank() []utils.RankItem {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	ret := []utils.RankItem{}
	if len(cache.rankNodes) > 0 {
		for _, item := range cache.rankNodes[len(cache.rankNodes)-1] {
			ret = append(ret, utils.RankItem{
				UserId: item.userId,
				Score:  item.score,
			})
		}
	}
	// TODO: add zero score users
	return ret
}

func (cache *TimelineRankCache) GetUserScore(userId uint64) uint64 {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	if len(cache.rankNodes) == 0 {
		return 0
	}
	for _, item := range cache.rankNodes[len(cache.rankNodes)-1] {
		if item.userId == userId {
			return item.score
		}
	}
	return 0
}

func (cache *TimelineRankCache) refreshRank(index int) {
	// TODO: step 1 re-calculate
	// TODO: step 2 sort
}

func (cache *TimelineRankCache) MutateUser(userId uint64, state bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.userState[userId] = state

	for index, frame := range cache.frames {
		if user, ok := frame.Users[userId]; ok {
			if len(user.Solved) != 0 {
				cache.refreshRank(index)
			}
		}
	}
}

func (cache *TimelineRankCache) MutateChallenge(challengeId uint64, state bool, dynamic bool, baseScore uint64) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.challengeState[challengeId] = state
	cache.challengeDynamic[challengeId] = dynamic
	cache.challengeScore[challengeId] = baseScore

	for index, frame := range cache.frames {
		if _, ok := frame.Challenges[challengeId]; ok {
			cache.refreshRank(index)
		}
	}
}

func (cache *TimelineRankCache) Submit(userId uint64, challengeId uint64, stamp time.Time) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	newFrame := TimeFrame{
		Challenges: map[uint64]*ChallengeFrame{},
		Users:      map[uint64]*UserFrame{},
		Stamp:      stamp,
	}

	if len(cache.frames) > 0 {
		lastFrame := cache.frames[len(cache.frames)-1]
		for k, v := range lastFrame.Challenges {
			newFrame.Challenges[k] = &ChallengeFrame{
				Solved: v.Solved,
			}
		}
		for k, user := range lastFrame.Users {
			item := &UserFrame{
				Solved: []uint64{},
				Blood:  map[uint64]int{},
			}
			for _, it := range user.Solved {
				item.Solved = append(newFrame.Users[k].Solved, it)
			}
			for bk, bv := range user.Blood {
				item.Blood[bk] = bv
			}
			newFrame.Users[k] = item
		}
	}

	if _, ok := newFrame.Challenges[challengeId]; !ok {
		newFrame.Challenges[challengeId] = &ChallengeFrame{
			Solved: 0,
		}
	}
	newFrame.Challenges[challengeId].Solved = newFrame.Challenges[challengeId].Solved + 1

	if _, ok := newFrame.Users[userId]; !ok {
		newFrame.Users[userId] = &UserFrame{
			Solved: []uint64{challengeId},
			Blood:  map[uint64]int{},
		}
	}
	solveIndex := newFrame.Challenges[challengeId].Solved - 1
	newFrame.Users[userId].Solved = append(newFrame.Users[userId].Solved, challengeId)
	if solveIndex < 3 {
		newFrame.Users[userId].Blood[challengeId] = int(solveIndex)
	}

	err := cache.fileCache.Append(&newFrame)
	if err != nil {
		return err
	}

	cache.frames = append(cache.frames, newFrame)
	cache.rankNodes = append(cache.rankNodes, ScoreItems{})
	return nil
}

func (cache *TimelineRankCache) GetCurrentScores() map[uint64]uint64 {
	ret := map[uint64]uint64{}
	for id, _ := range cache.challengeState {
		solved := uint64(0)
		if len(cache.frames) != 0 {
			if challenge, ok := cache.frames[len(cache.frames)-1].Challenges[id]; ok {
				solved = challenge.Solved
			}
		}
		ret[id] = cache.calculator.GetScore(cache.challengeScore[id], solved)
	}
	return ret
}

func (cache *TimelineRankCache) Load() error {
	data, err := cache.fileCache.Load()
	if err != nil {
		return err
	}
	cache.frames = data
	cache.rankNodes = make([]ScoreItems, len(cache.frames))

	cache.challengeState = map[uint64]bool{}
	cache.challengeDynamic = map[uint64]bool{}
	cache.challengeScore = map[uint64]uint64{}
	cache.userState = map[uint64]bool{}
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
		cache.MutateChallenge(challenge.ChallengeId, true, config.Score.Dynamic, score)
	}
	users := resolvers.FindAllUser()
	if users != nil {
		for _, user := range users {
			if user.Role != "admin" {
				cache.MutateUser(user.UserId, user.State != "disabled")
			}
		}
	}

	for i, _ := range cache.frames {
		cache.refreshRank(i)
	}

	return nil
}
