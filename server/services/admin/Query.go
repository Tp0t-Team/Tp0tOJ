package admin

import (
	"context"
	"github.com/kataras/go-sessions/v3"

	"server/services/database/resolvers"
	"server/services/types"
	"strconv"
)

type QueryResolver struct {
}

func (r *QueryResolver) AllUserInfos(ctx context.Context) (*types.AllUserInfoResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.AllUserInfoResult{Message: "forbidden"}, nil
	}
	result := types.AllUserInfoResult{Message: ""}
	users, err := resolvers.FindAllUser()
	if err != nil {
		return nil, err
	}
	//TODO: May need some method to cache
	for _, v := range users {
		result.AllUserInfos = append(result.AllUserInfos, types.UserInfo{UserId: strconv.FormatUint(v.UserId, 10), Name: v.Name, Avatar: v.MakeAvatarUrl(), Mail: v.Mail, JoinTime: v.JoinTime.String(), Score: strconv.FormatInt(v.Score, 10), Role: v.Role, State: v.State})
	}
	return &result, nil
}

func (r *QueryResolver) ChallengeConfigs(ctx context.Context) (*types.ChallengeConfigsResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.ChallengeConfigsResult{Message: "forbidden"}, nil
	}

}

func (r *QueryResolver) SubmitHistory(userId string, ctx context.Context) (*types.SubmitHistoryResult, error) {

}
