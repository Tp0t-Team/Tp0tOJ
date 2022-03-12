package services

import (
	"context"
	_ "embed"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/kataras/go-sessions/v3"
	"net/http"
	"server/services/admin"
	"server/services/user"
	"time"
)

type Resolver struct {
	admin.AdminMutationResolver
	admin.AdminQueryResolver
	user.MutationResolver
	user.QueryResolver
}

//go:embed schema.graphql
var schemaStr string

func init() {
	sessionManager := sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "gosessionid"
		Cookie: "session",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 1,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it
		DisableSubdomainPersistence: false,
		// want to be crazy safe? Take a look at the "securecookie" example folder.
	})

	schema := graphql.MustParseSchema(schemaStr, &Resolver{}, graphql.UseFieldResolvers())
	//http.Handle("/query", &relay.Handler{Schema: schema})
	graphqlHandle := &relay.Handler{Schema: schema}
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		graphqlHandle.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "session", session)))
	})
	http.HandleFunc("/writeup", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		isLogin := session.Get("isLogin")
		if isLogin == nil || !*isLogin.(*bool) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
			return
		}
		userId := *session.Get("userId").(*uint64)
		user.WriteUpHandle(w, r, userId)
	})
	http.HandleFunc("/wp", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		isLogin := session.Get("isLogin")
		isAdmin := session.Get("isAdmin")
		if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
			return
		}
		params := r.URL.Query()
		userId := params.Get("userId")
		if userId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}
		admin.DownloadWPByUserId(w, r, userId)
	})
	http.HandleFunc("/allwp", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		isLogin := session.Get("isLogin")
		isAdmin := session.Get("isAdmin")
		if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
			return
		}
		admin.DownloadAllWP(w, r)
	})
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		isLogin := session.Get("isLogin")
		isAdmin := session.Get("isAdmin")
		if isLogin == nil || !*isLogin.(*bool) || isAdmin == nil || !*isAdmin.(*bool) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(nil)
			return
		}
		admin.UploadImage(w, r)
	})
}
