package admin

import (
	"context"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/services/database"
	"server/services/types"
	"strconv"
)

type MutationResolver struct {
}

func (r *MutationResolver) BulletinPub(input types.BulletinPubInput, ctx context.Context) (*types.BulletinPubResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.BulletinPubResult{Message: "forbidden"}, nil
	}
	if input.CheckPass() {
		return &types.BulletinPubResult{Message: "not empty error"}, nil
	}
	err := database.AddBulletin(input.Title, input.Content, input.Topping)
	if err != nil {
		log.Println("Bulletin addition Error: ", err.Error())
		return &types.BulletinPubResult{Message: "Bulletin addition Error!"}, nil
	}
	return &types.BulletinPubResult{Message: ""}, nil

}
func (r *MutationResolver) UserInfoUpdate(input types.UserInfoUpdateInput, ctx context.Context) (*types.UserInfoUpdateResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	userId := session.Get("userId").(*uint64)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.UserInfoUpdateResult{Message: "forbidden"}, nil
	}
	if input.CheckPass() {
		return &types.UserInfoUpdateResult{Message: "user information check failed"}, nil
	}
	if userId != nil {
		checkResult, err := database.CheckAdminByUserId(*userId)
		if err != nil {
			return nil, err
		}
		inputUserId, err := strconv.ParseUint(input.UserId, 10, 64)
		if err != nil {
			log.Println("userId parse error", err)
			return nil, err
		}
		if checkResult && inputUserId == *userId && !(input.Role == "admin") {
			return &types.UserInfoUpdateResult{Message: "downgrade not permitted"}, nil
		}
		err = database.UpdateUserInfo(inputUserId, input.Name, input.Role, input.Mail, input.State)
		if err != nil {
			return nil, err
		}
		return &types.UserInfoUpdateResult{Message: ""}, nil

	}

}

func (r *MutationResolver) ChallengeMutate(input types.ChallengeMutateInput, ctx context.Context) (*types.ChallengeMutateResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.ChallengeMutateResult{Message: "forbidden"}, nil
	}
	if input.CheckPass() {
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}, nil
	}
	if input.ChallengeId == "" {
		//database.
		//TODO:
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}, nil
	}
	inputChallengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
	if err != nil {
		log.Println("challengeId parse error", err)
	}
}
func (r *MutationResolver) ChallengeRemove(id string) (*types.ChallengeRemoveResult, error) {
	return &types.ChallengeRemoveResult{Message: "Can't remove any challenge"}, nil
}
func (r *MutationResolver) WarmUp() (bool, error) {

}
