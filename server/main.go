package main

import (
	"fmt"
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
	err := utils.Cache.WarmUp()
	if err != nil {
		log.Panicln(err)
	}

	_, err = os.Stat(configure.WriteUpPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(configure.WriteUpPath, os.FileMode(0755))
			if err != nil {
				log.Panicln("writeup make dir error", err)
			}
		} else {
			log.Panicln("writeup dir create filed", err)
		}
	}

	_, crtErr := os.Stat("resources/https.crt")
	_, keyErr := os.Stat("resources/https.key")
	if crtErr == nil && keyErr == nil {
		if configure.Configure.Server.Port == 0 {
			configure.Configure.Server.Port = 443
		}
		portString := fmt.Sprintf(":%d", configure.Configure.Server.Port)
		log.Fatal(http.ListenAndServeTLS(portString, "resources/https.crt", "resources/https.key", nil))
	} else {
		if configure.Configure.Server.Port == 0 {
			configure.Configure.Server.Port = 80
		}
		portString := fmt.Sprintf(":%d", configure.Configure.Server.Port)
		log.Fatal(http.ListenAndServe(portString, nil))
	}
}
