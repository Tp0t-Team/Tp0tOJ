package main

import (
	//"github.com/gorilla/sessions"
	"log"
	_ "server/services"
	"server/utils"
)

func main() {
	// TODO: provide --prepare flags to check environment and install k3s server and other requirement

	config := new(utils.Config)
	utils.Parse(config)
	//http.Handle("")
	log.Println(config)
}
