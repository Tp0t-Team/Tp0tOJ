package user

import (
	"context"
	"encoding/json"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/entity"
	"server/services/database/resolvers"
	"server/services/kube"
	"server/services/types"
	"server/utils"
	"server/utils/kick"
	"strconv"
	"time"
)

var timeZone *time.Location

func init() {
	var err error
	timeZone, err = time.LoadLocation("UTC")
	if err != nil {
		log.Panicln(err)
	}
}

type QueryResolver struct {
}

func (r *QueryResolver) AllBulletin() *types.BulletinResult {
	bulletins := resolvers.GetAllBulletin()
	if bulletins == nil {
		return &types.BulletinResult{Message: "Get Bulletin Error!"}
	}
	result := types.BulletinResult{Message: "", Bulletins: []types.Bulletin{}}
	for _, bulletin := range bulletins {
		result.Bulletins = append(result.Bulletins, types.Bulletin{
			Title:       bulletin.Title,
			Content:     bulletin.Content,
			PublishTime: bulletin.PublishTime.In(timeZone).Format(time.RFC3339),
		})
	}
	return &result
}

func (r *QueryResolver) Rank(ctx context.Context) *types.RankResult {
	rank := utils.Cache.GetRank()
	ret := []types.RankResultDesc{}
	for _, item := range rank {
		user, err := resolvers.FindUser(item.UserId)
		if err != nil || user == nil {
			return &types.RankResult{Message: "get rank error"}
		}
		if user.State == "disabled" {
			continue
		}
		ret = append(ret, types.RankResultDesc{
			UserId: strconv.FormatUint(item.UserId, 10),
			Name:   user.Name,
			Avatar: user.MakeAvatarUrl(),
			Score:  int32(item.Score),
		})
	}
	return &types.RankResult{Message: "", RankResultDescs: ret}
}

func (r *QueryResolver) UserInfo(ctx context.Context, args struct{ UserId string }) *types.UserInfoResult {
	userId := args.UserId
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.UserInfoResult{Message: "unauthorized"}
	}
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.UserInfoResult{Message: "Get User Info Error!"}
	}
	currentUserId := *session.Get("userId").(*uint64)
	if !kick.KickGuard(currentUserId) {
		return &types.UserInfoResult{Message: "forbidden"}
	}
	var user *entity.User
	user, err = resolvers.FindUser(parsedUserId)
	if err != nil {
		log.Println(err)
		return &types.UserInfoResult{Message: "Get User Info Error!"}
	}
	if user == nil {
		return &types.UserInfoResult{Message: "No such user."}
	}
	result := types.UserInfoResult{
		Message: "",
		UserInfo: types.UserInfo{
			UserId:   strconv.FormatUint(user.UserId, 10),
			Name:     user.Name,
			Avatar:   user.MakeAvatarUrl(),
			Mail:     "",
			JoinTime: user.JoinTime.In(timeZone).Format(time.RFC3339),
			Score:    int32(utils.Cache.GetUserScore(user.UserId)),
			Role:     user.Role,
			State:    user.State,
			//Rank:     0, //
		}}
	if *session.Get("isAdmin").(*bool) || parsedUserId == currentUserId {
		result.UserInfo.Mail = user.Mail
	}
	return &result
}

func (r *QueryResolver) ChallengeInfos(ctx context.Context) *types.ChallengeInfosResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.ChallengeInfosResult{Message: "unauthorized"}
	}
	currentUserId := *session.Get("userId").(*uint64)
	if !kick.KickGuard(currentUserId) {
		return &types.ChallengeInfosResult{Message: "forbidden"}
	}
	if !resolvers.IsGameRunning(nil) && !*isAdmin.(*bool) {
		return &types.ChallengeInfosResult{Message: "game is not running now"}
	}
	challenges := resolvers.FindEnabledChallenges()
	if challenges == nil {
		return &types.ChallengeInfosResult{Message: "Get Challenge Info Error!"}
	}
	result := types.ChallengeInfosResult{
		Message:        "",
		ChallengeInfos: []types.ChallengeInfo{},
	}
	scoreSet := utils.Cache.GetCurrentScores()
	for _, challenge := range challenges {
		var bloodInfo []types.BloodInfo
		if challenge.FirstBloodId != nil {
			user, err := resolvers.FindUser(*challenge.FirstBloodId)
			if err != nil {
				return nil
			}
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(user.UserId, 10),
				Name:   user.Name,
				Avatar: user.MakeAvatarUrl(),
			})
		}
		if challenge.SecondBloodId != nil {
			user, err := resolvers.FindUser(*challenge.SecondBloodId)
			if err != nil {
				return nil
			}
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(user.UserId, 10),
				Name:   user.Name,
				Avatar: user.MakeAvatarUrl(),
			})
		}
		if challenge.ThirdBloodId != nil {
			user, err := resolvers.FindUser(*challenge.ThirdBloodId)
			if err != nil {
				return nil
			}
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(user.UserId, 10),
				Name:   user.Name,
				Avatar: user.MakeAvatarUrl(),
			})
		}

		correct := resolvers.CheckSubmitCorrectByUserIdAndChallengeId(currentUserId, challenge.ChallengeId, nil)
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.ChallengeInfosResult{Message: "Get Challenge Info Error!"}
		}

		var realScore int32
		if score, ok := scoreSet[challenge.ChallengeId]; ok {
			realScore = int32(score)
		} else {
			tmp, _ := strconv.ParseInt(config.Score.BaseScore, 10, 32)
			realScore = int32(tmp)
		}
		var solved int
		submits := resolvers.FindSubmitCorrectByChallengeId(challenge.ChallengeId)
		if submits == nil {
			solved = 0
		}
		solved = len(submits)
		item := types.ChallengeInfo{
			ChallengeId: strconv.FormatUint(challenge.ChallengeId, 10),
			Category:    config.Category,
			Name:        challenge.Name,
			Score:       realScore,
			Blood:       bloodInfo,
			SolvedNum:   int32(solved),
			Done:        correct,
		}
		result.ChallengeInfos = append(result.ChallengeInfos, item)
	}
	return &result
}

func (r *QueryResolver) WatchDescription(ctx context.Context, args struct{ ChallengeId string }) *types.WatchDescriptionResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.WatchDescriptionResult{Message: "unauthorized"}
	}
	currentUserId := *session.Get("userId").(*uint64)
	if !kick.KickGuard(currentUserId) {
		return &types.WatchDescriptionResult{Message: "forbidden"}
	}
	if !resolvers.IsGameRunning(nil) && !*isAdmin.(*bool) {
		return &types.WatchDescriptionResult{Message: "game is not running now"}
	}
	parsedChallengeId, err := strconv.ParseUint(args.ChallengeId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.WatchDescriptionResult{Message: "service error"}
	}
	err = resolvers.AllocSingleton(parsedChallengeId, currentUserId)
	if err != nil {
		resolvers.BehaviorWatchDescription(parsedChallengeId, currentUserId, time.Now(), nil)
		return &types.WatchDescriptionResult{Message: err.Error()}
	}
	challenge, err := resolvers.FindChallengeById(parsedChallengeId)
	if err != nil {
		log.Println(err)
		return &types.WatchDescriptionResult{Message: "service error"}
	}
	var config types.ChallengeConfig
	err = json.Unmarshal([]byte(challenge.Configuration), &config)
	if err != nil {
		log.Println(err)
		return &types.WatchDescriptionResult{Message: "Get Challenge Description Error!"}
	}
	replicaUrls := []string{}
	var alloc *entity.ReplicaAlloc
	alloc, err = resolvers.FindReplicaAllocByUserIdAndChallengeId(currentUserId, parsedChallengeId, nil)
	if err != nil {
		return &types.WatchDescriptionResult{Message: "Error"}
	}
	if alloc != nil {
		for _, node := range config.NodeConfig {
			replicaUrls = append(replicaUrls, kube.K8sServiceGetUrls(alloc.ReplicaId, node.Name)...)
		}
	}
	replicaUrls = append(replicaUrls, config.ExternalLink...)
	allocState := types.AllocatedUndone
	if alloc != nil {
		allocState = types.AllocatedDone
	} else {
		resolvers.AllocatingTableMtx.RLock()
		if _, ok := resolvers.AllocatingTable[currentUserId]; ok {
			if _, ok := resolvers.AllocatingTable[currentUserId][challenge.ChallengeId]; ok {
				allocState = types.AllocatedDoing
			}
		}
		resolvers.AllocatingTableMtx.RUnlock()
	}
	description := types.ChallengeDesc{
		ChallengeId:  strconv.FormatUint(parsedChallengeId, 10),
		Description:  config.Description,
		ExternalLink: replicaUrls,
		Manual:       !config.Singleton && len(config.NodeConfig) > 0,
		Allocated:    allocState,
	}
	resolvers.BehaviorWatchDescription(parsedChallengeId, currentUserId, time.Now(), nil)
	return &types.WatchDescriptionResult{Message: "", Description: description}
}
