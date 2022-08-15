package admin

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

type AdminQueryResolver struct {
}

func (r *AdminQueryResolver) AllUserInfos(ctx context.Context) *types.AllUserInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.AllUserInfoResult{Message: "forbidden or login timeout"}
	}
	var userInfos []types.UserInfo
	users := resolvers.FindAllUser()
	if users == nil {
		return &types.AllUserInfoResult{Message: "Get User Info Error!", AllUserInfos: userInfos}
	}
	//TODO: May need some method to cache
	for _, v := range users {
		userInfos = append(userInfos, types.UserInfo{UserId: strconv.FormatUint(v.UserId, 10), Name: v.Name, Avatar: v.MakeAvatarUrl(), Mail: v.Mail, JoinTime: v.JoinTime.String(), Score: int32(utils.Cache.GetUserScore(v.UserId)), Role: v.Role, State: v.State})
	}
	return &types.AllUserInfoResult{Message: "", AllUserInfos: userInfos}
}

func (r *AdminQueryResolver) ChallengeConfigs(ctx context.Context) *types.ChallengeConfigsResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.ChallengeConfigsResult{Message: "forbidden or login timeout"}
	}
	challenges := resolvers.FindAllChallenges()
	if challenges == nil {
		return &types.ChallengeConfigsResult{Message: "Challenge Config Error!"}
	}
	var challengeConfigs []types.ChallengeConfigAndState
	for _, challenge := range challenges {
		var config types.ChallengeConfig
		err := json.Unmarshal([]byte(challenge.Configuration), &config)
		if err != nil {
			log.Println(err)
			return &types.ChallengeConfigsResult{Message: "Challenge Config Error!"}
		}
		challengeConfigs = append(challengeConfigs, types.ChallengeConfigAndState{
			ChallengeId: strconv.FormatUint(challenge.ChallengeId, 10),
			Name:        challenge.Name,
			Config:      config,
			State:       challenge.State,
		})
	}
	return &types.ChallengeConfigsResult{Message: "", ChallengeConfigs: challengeConfigs}
}

func (r *AdminQueryResolver) SubmitHistory(ctx context.Context, args struct{ UserId string }) *types.SubmitHistoryResult {
	userId := args.UserId
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.SubmitHistoryResult{Message: "forbidden or login timeout"}
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
		submitInfo := types.SubmitInfo{SubmitTime: submit.SubmitTime.In(timeZone).Format(time.RFC3339), ChallengeName: submit.Challenge.Name}
		submitInfos = append(submitInfos, submitInfo)
	}
	return &types.SubmitHistoryResult{Message: "", SubmitInfos: submitInfos}

}

func (r *AdminQueryResolver) WriteUpInfos(ctx context.Context) *types.WriteUpInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.WriteUpInfoResult{Message: "forbidden or login timeout"}
	}
	ret := GetWriteUpInfos()
	return &types.WriteUpInfoResult{Message: "", Infos: ret}
}

func (r *AdminQueryResolver) ImageInfos(ctx context.Context) *types.ImageInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.ImageInfoResult{Message: "forbidden or login timeout"}
	}
	ret := kube.ImgStatus()
	if ret == nil {
		ret = []types.ImageInfo{}
	}
	return &types.ImageInfoResult{Message: "", Infos: ret}
}

func (r *AdminQueryResolver) ClusterInfo(ctx context.Context) *types.ClusterInfoResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.ClusterInfoResult{Message: "forbidden or login timeout"}
	}
	nodeInfos, replicaInfos := kube.K8sStatus()
	for _, infoItem := range replicaInfos {
		if _, ok := kube.DeletingReplicas[infoItem.Name]; ok {
			infoItem.Status = "deleting"
		}
	}
	return &types.ClusterInfoResult{
		Message:  "",
		Nodes:    nodeInfos,
		Replicas: replicaInfos,
	}
}

func (r *AdminMutationResolver) AllEvents(ctx context.Context) *types.AllEventResult {

	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
		return &types.AllEventResult{Message: "forbidden or login timeout"}
	}
	var events []entity.GameEvent
	events = resolvers.GetAllEvents()
	if events == nil {
		return &types.AllEventResult{Message: "get events error"}
	}
	var results []types.GameEvent

	for _, event := range events {
		result := types.GameEvent{
			EventId: strconv.FormatUint(event.EventId, 10),
			Time:    strconv.FormatInt(event.Time.Unix(), 10),
			Action:  int32(event.Action),
		}
		results = append(results, result)
	}

	return &types.AllEventResult{Message: "", AllEvents: results}
}
