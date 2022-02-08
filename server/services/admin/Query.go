package admin

import (
	"context"
	"server/services/types"
)

type QueryResolver struct {
}

func (r *QueryResolver) AllUserInfos(ctx context.Context) (*types.AllUserInfoResult, error) {

}

func (r *QueryResolver) ChallengeConfigs(ctx context.Context) (*types.ChallengeConfigsResult, error) {

}

func (r *QueryResolver) SubmitHistory(userId string, ctx context.Context) (*types.SubmitHistoryResult, error) {

}
