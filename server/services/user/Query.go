package user

import (
	"context"
	"github.com/kataras/go-sessions/v3"
	"server/entity"
	"server/services/database"
	"server/services/types"
	"strconv"
	"time"
)

type QueryResolver struct {
}

func (r *QueryResolver) AllBulletin() (*types.BulletinResult, error) {
	bulletins, err := database.GetAllBulletin()
	if err != nil {
		return &types.BulletinResult{Message: err.Error()}, nil
	}
	result := types.BulletinResult{Message: "", Bulletins: []types.Bulletin{}}
	for _, bulletin := range bulletins {
		result.Bulletins = append(result.Bulletins, types.Bulletin{
			Title:       bulletin.Title,
			Content:     bulletin.Content,
			PublishTime: bulletin.PublishTime.Format(time.RFC3339),
		})
	}
	return &result, nil
}

func (r *QueryResolver) Rank(ctx context.Context) (*types.RankResult, error) {
	// TODO:
}

func (r *QueryResolver) UserInfo(userId string, ctx context.Context) (*types.UserInfoResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.UserInfoResult{Message: "forbidden"}, nil
	}
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return &types.UserInfoResult{Message: err.Error()}, nil
	}
	currentUserId := *session.Get("userId").(*uint64)
	var user *entity.User
	user, err = database.FindUser(parsedUserId)
	if err != nil {
		return &types.UserInfoResult{Message: err.Error()}, nil
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
			Rank:     0, // TODO:
		}}
	if *session.Get("isAdmin").(*bool) || parsedUserId == currentUserId {
		result.UserInfo.Mail = user.Mail
	}
	return &result, nil
}

func (r *QueryResolver) ChallengeInfos(ctx context.Context) (*types.ChallengeInfosResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.ChallengeInfosResult{Message: "forbidden"}, nil
	}
	currentUserId := *session.Get("userId").(*uint64)
	challenges, err := database.FindEnabledChallenges()
	if err != nil {
		return &types.ChallengeInfosResult{Message: err.Error()}, nil
	}
	result := types.ChallengeInfosResult{
		Message:        "",
		ChallengeInfos: []types.ChallengeInfo{},
	}
	for _, challenge := range challenges {
		bloodInfo := []types.BloodInfo{}
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
		correct, _ := database.CheckSubmitCorrectByUserIdAndChallengeId(currentUserId, challenge.ChallengeId)
		item := types.ChallengeInfo{
			ChallengeId:  strconv.FormatUint(challenge.ChallengeId, 10),
			Category:     "",  // TODO:
			Name:         "",  // TODO:
			Score:        0,   // TODO:
			Description:  "",  // TODO:
			ExternalLink: nil, // TODO:
			Hint:         nil, // TODO:
			Blood:        bloodInfo,
			Done:         correct,
		}
		result.ChallengeInfos = append(result.ChallengeInfos, item)
	}
	return &result, nil
}
