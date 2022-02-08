package main

import (
	//"github.com/gorilla/sessions"
	"log"
	_ "server/services"
	"server/utils"
)

func main() {
	config := new(utils.Config)
	utils.Parse(config)
	//http.Handle("")
	log.Println(config)
}
