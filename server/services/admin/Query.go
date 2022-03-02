package admin

import (
	"context"
	"encoding/json"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/services/database/resolvers"
	"server/services/types"
	"strconv"
)

type QueryResolver struct {
}

func (r *QueryResolver) AllUserInfos(ctx context.Context) *types.AllUserInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.AllUserInfoResult{Message: "forbidden"}
	}
	var userInfos []types.UserInfo
	users := resolvers.FindAllUser()
	if users == nil {
		return &types.AllUserInfoResult{Message: "Get User Info Error!", AllUserInfos: userInfos}
	}
	//TODO: May need some method to cache
	for _, v := range users {
		userInfos = append(userInfos, types.UserInfo{UserId: strconv.FormatUint(v.UserId, 10), Name: v.Name, Avatar: v.MakeAvatarUrl(), Mail: v.Mail, JoinTime: v.JoinTime.String(), Score: strconv.FormatInt(v.Score, 10), Role: v.Role, State: v.State})
	}
	return &types.AllUserInfoResult{Message: "", AllUserInfos: userInfos}
}

func (r *QueryResolver) ChallengeConfigs(ctx context.Context) *types.ChallengeConfigsResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.ChallengeConfigsResult{Message: "forbidden"}
	}
	challenges := resolvers.FindAllChallenges()
	if challenges == nil {
		return &types.ChallengeConfigsResult{Message: "Challenge Config Error!"}
	}
	var challengeConfigs []types.ChallengeConfig
	for _, challenge := range challenges {
		log.Println(challenge.Configuration)
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.ChallengeConfigsResult{Message: "Challenge Config Error!"}
		}
		challengeConfigs = append(challengeConfigs, config)
	}
	return &types.ChallengeConfigsResult{Message: "", ChallengeConfigs: challengeConfigs}
}

func (r *QueryResolver) SubmitHistory(userId string, ctx context.Context) *types.SubmitHistoryResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.SubmitHistoryResult{Message: "forbidden"}
	}

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.SubmitHistoryResult{Message: "Submit History Error!"}
	}
	submits := resolvers.FindSubmitCorrectByUserId(id)
	if submits == nil {
		return &types.SubmitHistoryResult{Message: "Submit History Error!"}
	}
	var submitInfos []types.SubmitInfo
	for _, submit := range submits {
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(submit.Challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.SubmitHistoryResult{Message: "Submit History Error!"}
		}
		submitInfo := types.SubmitInfo{SubmitTime: submit.SubmitTime.String(), ChallengeName: config.Name}
		submitInfos = append(submitInfos, submitInfo)
	}
	return &types.SubmitHistoryResult{Message: "", SubmitInfos: submitInfos}

}
