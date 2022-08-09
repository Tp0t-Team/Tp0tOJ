package admin

import (
	"context"
	"github.com/kataras/go-sessions/v3"
	"log"
	"server/services/database/resolvers"
	"server/services/kube"
	"server/services/types"
	"server/utils"
	"server/utils/kick"
	"strconv"
	"strings"
)

type AdminMutationResolver struct {
}

func (r *AdminMutationResolver) BulletinPub(ctx context.Context, args struct{ Input types.BulletinPubInput }) *types.BulletinPubResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.BulletinPubResult{Message: "forbidden or login timeout"}
	}
	if !input.CheckPass() {
		return &types.BulletinPubResult{Message: "not empty error"}
	}
	ok := resolvers.AddBulletin(input.Title, input.Content, input.Topping)
	if !ok {
		return &types.BulletinPubResult{Message: "resolvers addition Error!"}
	}
	return &types.BulletinPubResult{Message: ""}

}
func (r *AdminMutationResolver) UserInfoUpdate(ctx context.Context, args struct{ Input types.UserInfoUpdateInput }) *types.UserInfoUpdateResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.UserInfoUpdateResult{Message: "forbidden or login timeout"}
	}
	userId := session.Get("userId").(*uint64)
	if !input.CheckPass() {
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
		if input.State == "disabled" {
			kick.BanUser(inputUserId)
		} else {
			kick.UnbanUser(inputUserId)
		}
		return &types.UserInfoUpdateResult{Message: ""}

	}
	return &types.UserInfoUpdateResult{Message: "user ID is nil"}
}

func (r *AdminMutationResolver) ChallengeMutate(ctx context.Context, args struct{ Input types.ChallengeMutateInput }) *types.ChallengeMutateResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.ChallengeMutateResult{Message: "forbidden or login timeout"}
	}
	if !input.CheckPass() {
		return &types.ChallengeMutateResult{Message: "Challenge format not available"}
	}
	if input.ChallengeId == "" {
		ok := resolvers.AddChallenge(input)
		if !ok {
			return &types.ChallengeMutateResult{Message: "Add Challenge Error!"}
		}
		return &types.ChallengeMutateResult{Message: ""}
	}

	ok := resolvers.UpdateChallenge(input)
	if !ok {
		return &types.ChallengeMutateResult{Message: "Update Challenge Error!"}
	}
	return &types.ChallengeMutateResult{Message: ""}
}

// ChallengeAction Handle 3 types of action : [ enable | disable | delete ]
func (r *AdminMutationResolver) ChallengeAction(ctx context.Context, args struct{ Input types.ChallengeActionInput }) *types.ChallengeActionResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.ChallengeActionResult{Message: "forbidden or login timeout"}
	}
	if !input.CheckPass() {
		return &types.ChallengeActionResult{Message: "action format error"}
	}
	// TODO: maybe need some binary flag to mark which challenge occurred error
	if input.Action == "enable" {
		var ok = true
		for _, id := range input.ChallengeIds {
			ok = ok && resolvers.EnableChallengeById(id)
		}
		if !ok {
			return &types.ChallengeActionResult{Message: "enable challenges occurred some error "}
		}
	}
	if input.Action == "disable" {
		var ok = true
		for _, id := range input.ChallengeIds {
			ok = ok && resolvers.DisableChallengeById(id)
		}
		if !ok {
			return &types.ChallengeActionResult{Message: "disable challenges occurred some error "}
		}
	}
	if input.Action == "delete" {
		var ok = true
		for _, id := range input.ChallengeIds {
			ok = ok && resolvers.DeleteChallenge(id)
		}
		if !ok {
			return &types.ChallengeActionResult{Message: "delete challenges occurred some error "}
		}
	}

	return &types.ChallengeActionResult{Message: ""}
}

func (r *AdminMutationResolver) WarmUp() (bool, error) {
	err := utils.Cache.WarmUp()
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *AdminMutationResolver) DeleteImage(ctx context.Context, args struct{ Input string }) *types.DeleteImageResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.DeleteImageResult{Message: "forbidden or login timeout"}
	}
	//can not delete image when some challenge is enabled
	ok := resolvers.CheckEnabledChallengesByImage(input)
	if !ok {
		return &types.DeleteImageResult{Message: "some challenge rely this image, delete failed"}
	}
	err := kube.ImgDelete(input)
	if err != nil {
		log.Println(err)
		return &types.DeleteImageResult{Message: "delete image error"}
	}
	return &types.DeleteImageResult{Message: ""}
}

func (r *AdminMutationResolver) DeleteReplica(ctx context.Context, args struct{ Input string }) *types.DeleteReplicaResult {
	replicaName := args.Input
	replicaId, err := strconv.ParseUint(strings.Split(replicaName, "-")[1], 10, 64)
	if err != nil {
		log.Println(err)
		return &types.DeleteReplicaResult{Message: "wrong replica name"}
	}
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.DeleteReplicaResult{Message: "forbidden or login timeout"}
	}
	kube.DeletingReplicas[replicaName] = nil
	var deleted = func(status bool) {
		delete(kube.DeletingReplicas, replicaName)
	}
	if resolvers.DeleteReplicaById(replicaId, deleted) {
		return &types.DeleteReplicaResult{Message: ""}
	} else {
		go func() {
			kube.TaskQ <- kube.Task{
				Tasks: []interface{}{
					&kube.DestroyTask{
						ReplicaId:     replicaId,
						ContainerName: replicaName[len(strings.Split(replicaName, "-")[1])+9:],
					},
				},
				CB: deleted,
			}
		}()
	}
	return &types.DeleteReplicaResult{Message: "delete replica error"}

}
