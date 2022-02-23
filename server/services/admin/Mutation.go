package admin

import (
	"context"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils"
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
	err := resolvers.AddBulletin(input.Title, input.Content, input.Topping)
	if err != nil {
		log.Println("resolvers addition Error: ", err.Error())
		return &types.BulletinPubResult{Message: "resolvers addition Error!"}, nil
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
		checkResult, err := resolvers.CheckAdminByUserId(*userId)
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
		err = resolvers.UpdateUserInfo(inputUserId, input.Name, input.Role, input.Mail, input.State)
		if err != nil {
			return nil, err
		}
		return &types.UserInfoUpdateResult{Message: ""}, nil

	}
	return &types.UserInfoUpdateResult{Message: "user ID is nil"}, nil
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
		err := resolvers.AddChallenge(input)
		if err != nil {
			return nil, err
		}
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}, nil
	}

	err := resolvers.UpdateChallenge(input)
	if err != nil {
		return nil, err
	}
	return &types.ChallengeMutateResult{Message: ""}, nil
}

func (r *MutationResolver) ChallengeRemove(id string) (*types.ChallengeRemoveResult, error) {
	//TODO: impl
	return &types.ChallengeRemoveResult{Message: "Can't remove any challenge"}, nil
}
func (r *MutationResolver) WarmUp() (bool, error) {
	err := utils.Cache.WarmUp()
	if err != nil {
		return false, nil
	}
	return true, nil
}
