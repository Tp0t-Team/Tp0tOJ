package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/domodwyer/mailyak/v3"
	"github.com/kataras/go-sessions/v3"
	"io"
	"log"
	"net/smtp"
	"regexp"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils/configure"
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
	_, err := io.WriteString(hash1, configure.Configure.Server.Salt+password)
	if err != nil {
		log.Panicln(err.Error())
	}
	hash2 := sha256.New()
	_, err = io.WriteString(hash2, configure.Configure.Server.Salt+fmt.Sprintf("%x", hash1.Sum(nil)))
	if err != nil {
		log.Panicln(err.Error())
	}
	return fmt.Sprintf("%x", hash2.Sum(nil))
}

func (r *MutationResolver) Register(ctx context.Context, args struct{ Input types.RegisterInput }) *types.RegisterResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin != nil && *isLogin.(*bool) {
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

func (r *MutationResolver) Login(ctx context.Context, args struct{ Input types.LoginInput }) *types.LoginResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin != nil && *isLogin.(*bool) {
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
	if user.State == "disabled" {
		return &types.LoginResult{Message: "You have been disabled!"}
	}
	if user.Password != passwordHash(input.Password) {
		return &types.LoginResult{Message: "failed"}
	}
	state := true
	session.Set("isLogin", &state)
	userId := user.UserId
	session.Set("userId", &userId)
	var adminState bool
	if user.Role == "admin" {
		adminState = true
	} else {
		adminState = false
	}
	session.Set("isAdmin", &adminState)
	var teamState bool
	if user.Role == "team" {
		teamState = true
	} else {
		teamState = false
	}
	session.Set("isTeam", teamState)
	return &types.LoginResult{Message: "", UserId: strconv.FormatUint(user.UserId, 10), Role: user.Role}
}

func (r *MutationResolver) Logout(ctx context.Context) *types.LogoutResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.LogoutResult{Message: "not login yet"}
	}
	session.Delete("isTeam")
	session.Delete("isAdmin")
	session.Delete("userId")
	var state = false
	session.Set("isLogin", &state)
	return &types.LogoutResult{Message: ""}
}

func sendMail(address string, subject string, content string) bool {
	auth := smtp.PlainAuth("", configure.Configure.Email.Username, configure.Configure.Email.Password, configure.Configure.Email.Host)
	mail := mailyak.New(configure.Configure.Email.Host+":25", auth)
	mail.To(address)
	mail.From(configure.Configure.Email.Username)
	mail.Subject(subject)
	mail.Plain().Set(content)
	err := mail.Send()
	if err != nil {
		return false
	}
	return true
}

func (r *MutationResolver) Forget(args struct{ Input string }) *types.ForgetResult {
	input := args.Input
	if input == "" {
		return &types.ForgetResult{Message: "empty"}
	}
	user, err := resolvers.FindUserByMail(input)
	if err != nil {
		log.Println(err)
		return &types.ForgetResult{Message: "error"}
	}
	if user == nil {
		return &types.ForgetResult{Message: "no such user"}
	}
	result := resolvers.AddResetToken(user.UserId)
	if result == nil {
		return &types.ForgetResult{Message: "failed"}
	}
	if !sendMail(input, "password reset", fmt.Sprintf("Please use the follow link to reset your password.\\n%s/reset?token=%s", configure.Configure.Server.Host, result.Token)) {
		return &types.ForgetResult{Message: "send mail failed"}
	}
	return &types.ForgetResult{Message: ""}
}

func (r *MutationResolver) Reset(ctx context.Context, args struct{ Input types.ResetInput }) *types.ResetResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin != nil && *isLogin.(*bool) {
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

func (r *MutationResolver) Submit(ctx context.Context, args struct{ Input types.SubmitInput }) *types.SubmitResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.SubmitResult{Message: "forbidden"}
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
	ok := resolvers.AddSubmit(userId, challengeId, input.Flag, time.Now(), !*isAdmin)
	if !ok {
		return &types.SubmitResult{Message: "Submit Service Error!"}
	}

	// TODO: fill blood
	//challenge, err := resolvers.FindChallengeById(challengeId)
	//if err != nil {
	//	return nil
	//}
	//if challenge.FirstBloodId == nil {
	//	challenge.FirstBloodId =
	//}
	return &types.SubmitResult{Message: ""}

}
