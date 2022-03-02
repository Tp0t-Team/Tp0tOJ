package user

import (
	"context"
	"encoding/json"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/entity"
	"server/services/database/resolvers"
	"server/services/types"
	"strconv"
	"time"
)

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
			PublishTime: bulletin.PublishTime.Format(time.RFC3339),
		})
	}
	return &result
}

func (r *QueryResolver) Rank(ctx context.Context) *types.RankResult {
	// TODO:
	return nil
}

func (r *QueryResolver) UserInfo(userId string, ctx context.Context) *types.UserInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.UserInfoResult{Message: "forbidden"}
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
			JoinTime: user.JoinTime.Format(time.RFC3339),
			Score:    "0", // TODO: strconv.FormatInt(user.Score, 10),
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
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.ChallengeInfosResult{Message: "forbidden"}
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
		if challenge.FirstBlood != nil {
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(challenge.FirstBlood.UserId, 10),
				Name:   challenge.FirstBlood.Name,
				Avatar: challenge.FirstBlood.MakeAvatarUrl(),
			})
		}
		if challenge.SecondBlood != nil {
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(challenge.SecondBlood.UserId, 10),
				Name:   challenge.SecondBlood.Name,
				Avatar: challenge.SecondBlood.MakeAvatarUrl(),
			})
		}
		if challenge.ThirdBlood != nil {
			bloodInfo = append(bloodInfo, types.BloodInfo{
				UserId: strconv.FormatUint(challenge.ThirdBlood.UserId, 10),
				Name:   challenge.ThirdBlood.Name,
				Avatar: challenge.ThirdBlood.MakeAvatarUrl(),
			})
		}
		correct := resolvers.CheckSubmitCorrectByUserIdAndChallengeId(currentUserId, challenge.ChallengeId)
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.ChallengeInfosResult{Message: "Get Challenge Info Error!"}
		}
		// TODO: add thr replica url to external links
		item := types.ChallengeInfo{
			ChallengeId:  strconv.FormatUint(challenge.ChallengeId, 10),
			Category:     config.Category,
			Name:         config.Name,
			Score:        config.Score.BaseScore,
			Description:  config.Description,
			ExternalLink: config.ExternalLink,
			Hint:         config.Hint,
			Blood:        bloodInfo,
			Done:         correct,
		}
		result.ChallengeInfos = append(result.ChallengeInfos, item)
	}
	return &result
}
