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

func (r *MutationResolver) BulletinPub(input types.BulletinPubInput, ctx context.Context) *types.BulletinPubResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.BulletinPubResult{Message: "forbidden"}
	}
	if input.CheckPass() {
		return &types.BulletinPubResult{Message: "not empty error"}
	}
	ok := resolvers.AddBulletin(input.Title, input.Content, input.Topping)
	if !ok {
		return &types.BulletinPubResult{Message: "resolvers addition Error!"}
	}
	return &types.BulletinPubResult{Message: ""}

}
func (r *MutationResolver) UserInfoUpdate(input types.UserInfoUpdateInput, ctx context.Context) *types.UserInfoUpdateResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	userId := session.Get("userId").(*uint64)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.UserInfoUpdateResult{Message: "forbidden"}
	}
	if input.CheckPass() {
		return &types.UserInfoUpdateResult{Message: "user information check failed"}
	}
	if userId != nil {
		checkResult := resolvers.CheckAdminByUserId(*userId)
		inputUserId, err := strconv.ParseUint(input.UserId, 10, 64)
		if err != nil {
			log.Println("userId parse error", err)
			return &types.UserInfoUpdateResult{Message: "Update Error!"}
		}
		if checkResult && inputUserId == *userId && !(input.Role == "admin") {
			return &types.UserInfoUpdateResult{Message: "downgrade not permitted"}
		}
		ok := resolvers.UpdateUserInfo(inputUserId, input.Name, input.Role, input.Mail, input.State)
		if !ok {
			return &types.UserInfoUpdateResult{Message: "Update Error!"}
		}
		return &types.UserInfoUpdateResult{Message: ""}

	}
	return &types.UserInfoUpdateResult{Message: "user ID is nil"}
}

func (r *MutationResolver) ChallengeMutate(input types.ChallengeMutateInput, ctx context.Context) *types.ChallengeMutateResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	isAdmin := session.Get("isAdmin").(*bool)
	if isLogin == nil || !*isLogin || isAdmin == nil || !*isAdmin {
		return &types.ChallengeMutateResult{Message: "forbidden"}
	}
	if input.CheckPass() {
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}
	}
	if input.ChallengeId == "" {
		ok := resolvers.AddChallenge(input)
		if !ok {
			return &types.ChallengeMutateResult{Message: "Add Challenge Error!"}
		}
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}
	}

	ok := resolvers.UpdateChallenge(input)
	if !ok {
		return &types.ChallengeMutateResult{Message: "Add Challenge Error!"}
	}
	return &types.ChallengeMutateResult{Message: ""}
}

func (r *MutationResolver) ChallengeRemove(id string) *types.ChallengeRemoveResult {
	ok := resolvers.RemoveChallenge(id)
	if !ok {
		return &types.ChallengeRemoveResult{Message: "Remove challenge failed"}
	}
	_, err := r.WarmUp()
	if err != nil {
		return nil
	}
	return &types.ChallengeRemoveResult{Message: ""}
}

func (r *MutationResolver) WarmUp() (bool, error) {
	err := utils.Cache.WarmUp()
	if err != nil {
		return false, nil
	}
	return true, nil
}
