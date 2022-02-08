package services

import (
	"context"
	_ "embed"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/kataras/go-sessions/v3"
	"log"
	"net/http"
	"server/services/user"
	"time"
)

type Resolver struct {
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
		Cookie: "mysessionid",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 2,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it
		DisableSubdomainPersistence: false,
		// want to be crazy safe? Take a look at the "securecookie" example folder.
	})

	schema := graphql.MustParseSchema(schemaStr, &Resolver{})
	//http.Handle("/query", &relay.Handler{Schema: schema})
	graphqlHandle := &relay.Handler{Schema: schema}
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.Start(w, r)
		graphqlHandle.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "session", session)))
	})
	log.Fatal(http.ListenAndServe(":8888", nil))
}
