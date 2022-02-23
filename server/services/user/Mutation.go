package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/kataras/go-sessions/v3"
	"io"
	"log"
	"regexp"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils"
)

type MutationResolver struct {
}

var MailPattern *regexp.Regexp

func init() {
	MailPattern, _ = regexp.Compile("^[_A-Za-z0-9-+]+(.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(.[A-Za-z0-9]+)*(.[A-Za-z]{2,})$")
}

func passwordHash(password string) string {
	hash1 := sha256.New()
	_, err := io.WriteString(hash1, utils.Configure.Server.Salt+password)
	if err != nil {
		log.Panicln(err.Error())
	}
	hash2 := sha256.New()
	_, err = io.WriteString(hash2, utils.Configure.Server.Salt+fmt.Sprintf("%x", hash1.Sum(nil)))
	if err != nil {
		log.Panicln(err.Error())
	}
	return fmt.Sprintf("%x", hash2.Sum(nil))
}

func (r *MutationResolver) Register(input types.RegisterInput, ctx context.Context) (*types.RegisterResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.RegisterResult{Message: "already login cannot register"}, nil
	}
	if !input.CheckPass() {
		return &types.RegisterResult{Message: "invalid information"}, nil
	}
	if !MailPattern.MatchString(input.Mail) {
		return &types.RegisterResult{Message: "invalid mail"}, nil
	}
	err := resolvers.AddUser(input.Name, passwordHash(input.Password), input.Mail, "member", "normal")
	if err != nil {
		return &types.RegisterResult{Message: err.Error()}, nil
	}
	return &types.RegisterResult{Message: ""}, nil
}

func (r *MutationResolver) Login(input types.LoginInput, ctx context.Context) (*types.LoginResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.LoginResult{Message: "already login"}, nil
	}
	if !input.CheckPass() {
		return &types.LoginResult{Message: "not empty error"}, nil
	}
	user, err := resolvers.FindUserByMail(input.Mail)
	if err != nil {
		return &types.LoginResult{Message: err.Error()}, nil
	}
	if user.Password != passwordHash(input.Password) {
		return &types.LoginResult{Message: "failed"}, nil
	}
	state := true
	session.Set("isLogin", &state)
	userId := user.UserId
	session.Set("userId", userId)
	var adminState bool
	if user.Role == "admin" {
		adminState = true
	} else {
		adminState = false
	}
	session.Set("isAdmin", adminState)
	var teamState bool
	if user.Role == "team" {
		teamState = true
	} else {
		teamState = false
	}
	session.Set("isTeam", teamState)
	return &types.LoginResult{Message: ""}, nil
}

func (r *MutationResolver) Logout(ctx context.Context) (*types.LogoutResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.LogoutResult{Message: "not login yet"}, nil
	}
	session.Delete("isTeam")
	session.Delete("isAdmin")
	session.Delete("userId")
	var state = false
	session.Set("isLogin", &state)
	return &types.LogoutResult{Message: ""}, nil
}

func (r *MutationResolver) Forget(input string) (*types.ForgetResult, error) {
	// TODO:
	return nil, nil
}

func (r *MutationResolver) Reset(input types.ResetInput, ctx context.Context) (*types.ResetResult, error) {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.ResetResult{Message: "already login"}, nil
	}
	err := resolvers.ResetPassword(input.Token, passwordHash(input.Password))
	if err != nil {
		return &types.ResetResult{Message: err.Error()}, nil
	}
	return &types.ResetResult{Message: ""}, nil
}

func (r *MutationResolver) Submit(input types.SubmitInput) (*types.SubmitResult, error) {
	// TODO:
	return nil, nil
}
