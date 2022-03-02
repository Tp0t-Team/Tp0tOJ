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
	"strconv"
	"time"
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

func (r *MutationResolver) Register(input types.RegisterInput, ctx context.Context) *types.RegisterResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.RegisterResult{Message: "already login cannot register"}
	}
	if !input.CheckPass() {
		return &types.RegisterResult{Message: "invalid information"}
	}
	if !MailPattern.MatchString(input.Mail) {
		return &types.RegisterResult{Message: "invalid mail"}
	}
	ok := resolvers.AddUser(input.Name, passwordHash(input.Password), input.Mail, "member", "normal")
	if !ok {
		return &types.RegisterResult{Message: "Register Error!"}
	}
	return &types.RegisterResult{Message: ""}
}

func (r *MutationResolver) Login(input types.LoginInput, ctx context.Context) *types.LoginResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.LoginResult{Message: "already login"}
	}
	if !input.CheckPass() {
		return &types.LoginResult{Message: "not empty error"}
	}
	user, err := resolvers.FindUserByMail(input.Mail)
	if err != nil {
		log.Println(err)
		return &types.LoginResult{Message: "Login Service Error!"}
	}
	if user == nil {
		return &types.LoginResult{Message: "No such user."}
	}
	if user.Password != passwordHash(input.Password) {
		return &types.LoginResult{Message: "failed"}
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
	return &types.LoginResult{Message: ""}
}

func (r *MutationResolver) Logout(ctx context.Context) *types.LogoutResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin == nil || !*isLogin {
		return &types.LogoutResult{Message: "not login yet"}
	}
	session.Delete("isTeam")
	session.Delete("isAdmin")
	session.Delete("userId")
	var state = false
	session.Set("isLogin", &state)
	return &types.LogoutResult{Message: ""}
}

func (r *MutationResolver) Forget(input string) *types.ForgetResult {
	// TODO:
	return nil
}

func (r *MutationResolver) Reset(input types.ResetInput, ctx context.Context) *types.ResetResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.ResetResult{Message: "already login"}
	}
	if !input.CheckPass() {
		return &types.ResetResult{Message: "not empty error"}
	}
	ok := resolvers.ResetPassword(input.Token, passwordHash(input.Password))
	if !ok {
		return &types.ResetResult{Message: "Password Reset Service Error!"}
	}
	return &types.ResetResult{Message: ""}
}

func (r *MutationResolver) Submit(input types.SubmitInput, ctx context.Context) *types.SubmitResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin").(*bool)
	if isLogin != nil && *isLogin {
		return &types.SubmitResult{Message: "already login"}
	}
	if !input.CheckPass() {
		return &types.SubmitResult{Message: "not empty error"}
	}
	userId := *session.Get("userId").(*uint64)
	challengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.SubmitResult{Message: "Submit Service Error!"}
	}
	ok := resolvers.AddSubmit(userId, challengeId, input.Flag, time.Now())
	if !ok {
		return &types.SubmitResult{Message: "Submit Service Error!"}
	}
	return &types.SubmitResult{Message: ""}
}
