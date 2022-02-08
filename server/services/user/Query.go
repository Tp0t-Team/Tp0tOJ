package user

import (
	"context"
	"server/services/types"
)

type QueryResolver struct {
}

func (r *QueryResolver) AllBulletin() (*types.BulletinResult, error) {

}

func (r *QueryResolver) Rank(ctx context.Context) (*types.RankResult, error) {

}

func (r *QueryResolver) UserInfo(userId string, ctx context.Context) (*types.UserInfoResult, error) {

}

func (r *QueryResolver) ChallengeInfos(ctx context.Context) (*types.ChallengeInfosResult, error) {

}
