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
	"server/utils/configure"
	"server/utils/email"
	"server/utils/kick"
	"strconv"
	"strings"
	"sync"
	"time"
)

type resetTimer struct {
	clock map[uint64]*time.Timer
	lock  sync.RWMutex
}

func (t *resetTimer) NewTimer(userId uint64) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.clock[userId]; ok {
		return false
	}
	t.clock[userId] = time.NewTimer(5 * time.Minute)
	go func() {
		t.lock.RLock()
		clock := t.clock[userId]
		t.lock.RUnlock()
		select {
		case <-clock.C:
			clock.Stop()
			t.lock.Lock()
			if clock == t.clock[userId] {
				delete(t.clock, userId)
			}
			t.lock.Unlock()
			break
		case <-time.After(10 * time.Minute):
			log.Println("some thing error at resetTimer, timeout")
			break
		}
	}()
	return true

}

var ResetTimer = resetTimer{clock: map[uint64]*time.Timer{}}

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
	_, err = io.WriteString(hash2, configure.Configure.Server.Salt+fmt.Sprintf("%02x", hash1.Sum(nil)))
	if err != nil {
		log.Panicln(err.Error())
	}
	return fmt.Sprintf("%02x", hash2.Sum(nil))
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
	resolvers.BehaviorLogin(userId, ctx.Value("ip").(string), time.Now(), nil)
	return &types.LoginResult{Message: "", UserId: strconv.FormatUint(user.UserId, 10), Role: user.Role}
}

func (r *MutationResolver) Logout(ctx context.Context) *types.LogoutResult {
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.LogoutResult{Message: "unauthorized"}
	}
	session.Delete("isTeam")
	session.Delete("isAdmin")
	session.Delete("userId")
	var state = false
	session.Set("isLogin", &state)
	return &types.LogoutResult{Message: ""}
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
	if user.State == "disabled" {
		return &types.ForgetResult{Message: "You have been disabled!"}
	}
	ok := ResetTimer.NewTimer(user.UserId)
	if !ok {
		return &types.ForgetResult{Message: "you have recently reset password"}
	}
	result := resolvers.AddResetToken(user.UserId)
	if result == nil {
		return &types.ForgetResult{Message: "failed"}
	}
	content, err := email.RenderResetEmail(email.ResetInfo{
		Username: user.Name,
		Url:      fmt.Sprintf("%s:%s/reset?token=%s", configure.Configure.Server.Host, strconv.Itoa(configure.Configure.Server.Port), result.Token),
	})
	if err != nil {
		log.Println(err)
		return &types.ForgetResult{Message: "send mail failed"}
	}
	err = email.SendMail(input, "[Tp0tOJ] Reset your password", content)
	if err != nil {
		log.Println(err)
		return &types.ForgetResult{Message: "send mail failed"}
	}
	log.Println("send successfully ... ")
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
		return &types.SubmitResult{Message: "unauthorized"}
	}
	if !resolvers.IsGameRunning(nil) && !*isAdmin.(*bool) {
		return &types.SubmitResult{Message: "game is not running now"}
	}
	userId := *session.Get("userId").(*uint64)
	if !kick.KickGuard(userId) {
		return &types.SubmitResult{Message: "forbidden"}
	}
	if !input.CheckPass() {
		return &types.SubmitResult{Message: "not empty error"}
	}
	challengeId, err := strconv.ParseUint(input.ChallengeId, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.SubmitResult{Message: "Submit Service Error!"}
	}
	submitTime := time.Now()
	resolvers.BehaviorSubmit(challengeId, userId, input.Flag, submitTime, nil)
	ok, correct := resolvers.AddSubmit(userId, challengeId, input.Flag, submitTime, !*isAdmin.(*bool))
	if !ok {
		return &types.SubmitResult{Message: "Submit Service Error!"}
	}
	return &types.SubmitResult{Message: "", Correct: correct}

}

func (r *MutationResolver) StartReplica(ctx context.Context, args struct{ Input string }) *types.StartReplicaResult {
	input := args.Input
	session := ctx.Value("session").(*sessions.Session)
	isLogin := session.Get("isLogin")
	isAdmin := session.Get("isAdmin")
	if isLogin == nil || !*isLogin.(*bool) {
		return &types.StartReplicaResult{Message: "unauthorized"}
	}
	userId := *session.Get("userId").(*uint64)
	if !kick.KickGuard(userId) {
		return &types.StartReplicaResult{Message: "forbidden"}
	}
	if !resolvers.IsGameRunning(nil) && !*isAdmin.(*bool) {
		return &types.StartReplicaResult{Message: "game is not running now"}
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return &types.StartReplicaResult{Message: "not empty error"}
	}
	challengeId, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		log.Println(err)
		return &types.StartReplicaResult{Message: "start failed"}
	}
	//ok := ReplicaTimer.NewTimer(userId)
	//if !ok {
	//	return &types.StartReplicaResult{Message: "timer start failed"}
	//}
	// do start replica (close other not-singleton replica)
	if !resolvers.StartReplicaForUser(userId, challengeId) {
		return &types.StartReplicaResult{Message: "start failed"}
	}
	return &types.StartReplicaResult{Message: ""}
}
