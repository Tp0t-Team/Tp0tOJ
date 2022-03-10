package main

import (
	//"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	_ "server/services"
	"server/utils"
	"server/utils/calculator"
	"server/utils/configure"
	_ "server/utils/configure"
	_ "server/utils/rank"
)

func main() {
	// TODO: provide --prepare flags to check environment and install k3s server and other requirement
	utils.Cache.SetCalculator(&calculator.BasicScoreCalculator{})

	_, err := os.Stat(configure.WriteUpPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(configure.WriteUpPath, 600)
			if err != nil {
				log.Panicln("writeup make dir error", err)
			}
		} else {
			log.Panicln("writeup dir create filed", err)
		}
	}

	log.Fatal(http.ListenAndServe(":8888", nil))
}
