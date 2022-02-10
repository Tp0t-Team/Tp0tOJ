package admin

import (
	"context"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/services/database"
	"server/services/types"
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
	//TODO: check pass
	err := database.AddBulletin(input.Title, input.Content, input.Topping)
	if err != nil {
		log.Println("Bulletin addition Error: ", err.Error())
		return &types.BulletinPubResult{Message: "Bulletin addition Error!"}, nil
	}
	return &types.BulletinPubResult{Message: ""}, nil

}
func (r *MutationResolver) UserInfoUpdate(input types.ChallengeMutateInput, ctx context.Context) (*types.UserInfoUpdateResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.UserInfoUpdateResult{Message: "forbidden"}, nil
	}
}

func (r *MutationResolver) ChallengeMutate(input types.ChallengeMutateInput, ctx context.Context) (*types.ChallengeMutateResult, error) {

}
func (r *MutationResolver) ChallengeRemove(id string) (*types.ChallengeRemoveResult, error) {

}
func (r *MutationResolver) WarmUp() (bool, error) {

}
