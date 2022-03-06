package main

import (
	//"github.com/gorilla/sessions"
	"log"
	"net/http"
	_ "server/services"
)

func main() {
	// TODO: provide --prepare flags to check environment and install k3s server and other requirement

	log.Fatal(http.ListenAndServe(":8888", nil))
}
