package main

import (
	"fmt"
	//"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	_ "server/services"
	"server/services/database"
	"server/services/database/resolvers"
	"server/utils"
	"server/utils/calculator"
	"server/utils/configure"
	_ "server/utils/configure"
	_ "server/utils/rank"
)

func Redirect(w http.ResponseWriter, req *http.Request) {
	url := *req.URL
	url.Scheme = "https"
	target := url.String()
	log.Println(target)
	http.Redirect(w, req, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}

func main() {
	// setup database connection
	database.InitDB(configure.Configure.Database.Dsn)
	resolvers.InitDB(database.DataBase)

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
		go func() {
			err := http.ListenAndServe(portString, http.HandlerFunc(Redirect))
			if err != nil {
				log.Println(err)
			}
		}()
		log.Fatal(http.ListenAndServeTLS(portString, "resources/https.crt", "resources/https.key", nil))
	} else {
		if configure.Configure.Server.Port == 0 {
			configure.Configure.Server.Port = 80
		}
		portString := fmt.Sprintf(":%d", configure.Configure.Server.Port)
		log.Fatal(http.ListenAndServe(portString, nil))
	}
}
