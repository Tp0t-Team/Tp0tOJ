package services

import (
	_ "embed"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"log"
	"net/http"
	"server/services/user"
)

type Resolver struct {
	user.MutationResolver
}

//go:embed test.graphql
var schemaStr string

func init() {
	schema := graphql.MustParseSchema(schemaStr, &Resolver{})
	http.Handle("/query", &relay.Handler{Schema: schema})
	log.Fatal(http.ListenAndServe(":8888", nil))
}
