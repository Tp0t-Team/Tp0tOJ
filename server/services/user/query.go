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
		return &types.UserInfoResult{Message: "forbidden or login timeout"}
	}
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.UserInfoResult{Message: "Get User Info Error!"}
	}
	currentUserId := *session.Get("userId").(*uint64)
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
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.ChallengeInfosResult{Message: "forbidden or login timeout"}
	}
	currentUserId := *session.Get("userId").(*uint64)
	challenges := resolvers.FindEnabledChallenges()
	if challenges == nil {
		return &types.ChallengeInfosResult{Message: "Get Challenge Info Error!"}
	}
	result := types.ChallengeInfosResult{
		Message:        "",
		ChallengeInfos: []types.ChallengeInfo{},
	}
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

		correct := resolvers.CheckSubmitCorrectByUserIdAndChallengeId(currentUserId, challenge.ChallengeId)
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.ChallengeInfosResult{Message: "Get Challenge Info Error!"}
		}
		replicaUrls := []string{}
		var alloc *entity.ReplicaAlloc
		alloc, err = resolvers.FindReplicaAllocByUserIdAndChallengeId(currentUserId, challenge.ChallengeId, nil)
		if err != nil {
			return &types.ChallengeInfosResult{Message: "Error"}
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
		item := types.ChallengeInfo{
			ChallengeId:  strconv.FormatUint(challenge.ChallengeId, 10),
			Category:     config.Category,
			Name:         challenge.Name,
			Score:        config.Score.BaseScore,
			Description:  config.Description,
			ExternalLink: replicaUrls,
			Blood:        bloodInfo,
			Done:         correct,
			Manual:       !config.Singleton && len(config.NodeConfig) > 0, // TODO: maybe need more conditions
			Allocated:    allocState,
		}
		result.ChallengeInfos = append(result.ChallengeInfos, item)
	}
	return &result
}
