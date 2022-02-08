package user

import (
	"context"
	"server/services/types"
)

type MutationResolver struct {
}

func (r *MutationResolver) Register(input types.RegisterInput, ctx context.Context) (*types.RegisterResult, error) {

}
func (r *MutationResolver) Login(input types.LoginInput, ctx context.Context) (*types.LoginResult, error) {

}

func (r *MutationResolver) Logout(ctx context.Context) (*types.LogoutResult, error) {

}
func (r *MutationResolver) Forget(input string) (*types.ForgetResult, error) {

}
func (r *MutationResolver) Reset(input types.ResetInput) (*types.ResetResult, error) {

}

func (r *MutationResolver) Submit(input types.SubmitInput) (*types.SubmitResult, error) {

}
