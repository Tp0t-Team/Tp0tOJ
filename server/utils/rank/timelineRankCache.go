package rank

import (
	"encoding/json"
	"errors"
	"github.com/ugorji/go/codec"
	"io"
	"log"
	"os"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils"
	"sort"
	"strconv"
	"sync"
	"time"
)

type ChallengeFrame struct {
	Solved uint64
}

type UserFrame struct {
	Solved    []uint64
	Blood     map[uint64]int
	LastSolve time.Time
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

	fileCache  *FileCache
	chartCache *utils.ChartData
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
	utils.Cache = &TimelineRankCache{
		calculator:       nil,
		mutex:            sync.RWMutex{},
		rankNodes:        nil,
		frames:           nil,
		challengeState:   nil,
		challengeDynamic: nil,
		challengeScore:   nil,
		userState:        nil,
		fileCache:        nil,
		chartCache:       nil,
	}
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
	rank := []*ScoreItem{}
	frame := &cache.frames[index]

	bannedCount := map[uint64]uint64{}
	for id, user := range frame.Users {
		if len(user.Solved) == 0 {
			continue
		}
		if state, ok := cache.userState[id]; ok && state {
			continue
		}

		for _, challenge := range user.Solved {
			if _, ok := bannedCount[challenge]; !ok {
				bannedCount[challenge] = 0
			}
			bannedCount[challenge] = bannedCount[challenge] + 1
		}
	}

	realScore := map[uint64]uint64{}
	for id, challenge := range frame.Challenges {
		if state, ok := cache.challengeState[id]; ok && state {
			bc := uint64(0)
			if realBan, ok := bannedCount[id]; ok {
				bc = realBan
			}
			realScore[id] = cache.calculator.GetScore(cache.challengeScore[id], challenge.Solved-bc, cache.challengeDynamic[id])
		}
	}

	for id, user := range frame.Users {
		score := uint64(0)
		for _, challenge := range user.Solved {
			if _, ok := realScore[challenge]; !ok {
				continue
			}
			if index, ok := user.Blood[challenge]; ok {
				score += cache.calculator.GetIncrementScore(realScore[challenge], uint64(index))
			} else {
				score += realScore[challenge]
			}
		}
		rank = append(rank, &ScoreItem{
			userId: id,
			score:  score,
			stamp:  user.LastSolve,
		})
	}
	// TODO: step 2 sort
	sort.Sort(RankSorter{rank})

	cache.rankNodes[index] = rank
}

// MutateUser @state: true -> user is valid , false -> user is invalid
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

	cache.chartCache = nil
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

	cache.chartCache = nil
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
				item.Solved = append(item.Solved, it)
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
			Solved: []uint64{},
			Blood:  map[uint64]int{},
		}
	}
	solveIndex := newFrame.Challenges[challengeId].Solved - 1
	newFrame.Users[userId].Solved = append(newFrame.Users[userId].Solved, challengeId)
	if solveIndex < 3 {
		newFrame.Users[userId].Blood[challengeId] = int(solveIndex)
	}
	newFrame.Users[userId].LastSolve = stamp

	err := cache.fileCache.Append(&newFrame)
	if err != nil {
		return err
	}

	cache.frames = append(cache.frames, newFrame)
	cache.rankNodes = append(cache.rankNodes, ScoreItems{})
	cache.refreshRank(len(cache.frames) - 1)

	cache.chartCache = nil

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
		ret[id] = cache.calculator.GetScore(cache.challengeScore[id], solved, cache.challengeDynamic[id])
	}
	return ret
}

func (cache *TimelineRankCache) Load(filename string) error {
	var err error
	cache.fileCache, err = NewFileCache(filename)
	if err != nil {
		return err
	}
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

func (cache *TimelineRankCache) Chart(topN uint64) *utils.ChartData {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	if cache.chartCache != nil {
		return cache.chartCache
	}

	ret := &utils.ChartData{
		X: []int64{},
		Y: []*utils.ChartCurve{},
	}
	if len(cache.rankNodes) == 0 {
		cache.chartCache = ret
		return ret
	}

	for index, user := range cache.rankNodes[len(cache.rankNodes)-1] {
		if topN == 0 || uint64(index) < topN {
			userEntity, err := resolvers.FindUser(user.userId)
			if err != nil {
				log.Panicln(err)
			}
			ret.Y = append(ret.Y, &utils.ChartCurve{
				Id:    user.userId,
				Name:  userEntity.Name,
				Score: []uint64{},
			})
		}
	}

	for index, node := range cache.rankNodes {
		ret.X = append(ret.X, cache.frames[index].Stamp.UnixMilli())
		userScore := map[uint64]uint64{}
		for _, item := range node {
			userScore[item.userId] = item.score
		}

		for _, curve := range ret.Y {
			if score, ok := userScore[curve.Id]; ok {
				curve.Score = append(curve.Score, score)
			} else {
				curve.Score = append(curve.Score, 0)
			}
		}
	}

	cache.chartCache = ret
	return ret
}
