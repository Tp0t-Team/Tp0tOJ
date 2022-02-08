package main

import (
	"log"
	"server/utils"
)

func main() {
	config := new(utils.Config)
	utils.Parse(config)
	log.Println(config)
}
