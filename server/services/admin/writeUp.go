package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/services/database/resolvers"
	"strconv"
)

const writeUpPath string = "./writeup"

type WriteUpInfoResult struct {
	Name   string
	Mail   string
	Solved int //solved challenge number
}

func DownloadAll() {

}
func DownloadWPByUserId() {

}

func WriteUpInfos(w http.ResponseWriter, req *http.Request, userId uint64) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		w.Write(nil)
		return
	} //TODO: need login and role check

	_, err := os.Stat(writeUpPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(writeUpPath, 600)
			if err != nil {
				log.Panicln("writeup make dir error", err)
				return
			}
		} else {
			log.Panicln("writeup dir create filed", err)
			return
		}
	}
	files, err := filepath.Glob("./writeup")
	if err != nil {
		log.Panicln("writeup file", err)
		return
	}
	var writeUpInfos []WriteUpInfoResult
	for _, file := range files {
		userId, err := strconv.ParseUint(file, 10, 64)
		if err != nil {
			log.Panicln("writeup parse file error", err)
			return
		}
		user, err := resolvers.FindUser(userId)
		if err != nil {
			log.Panicln("writeup can't find user", err)
			return
		}
		var solved int
		submits := resolvers.FindSubmitCorrectByUserId(userId)
		if submits == nil {
			solved = 0
		}
		solved = len(submits)
		writeUpInfos = append(writeUpInfos, WriteUpInfoResult{
			Name:   user.Name,
			Mail:   user.Mail,
			Solved: solved,
		})
		infosJson, err := json.Marshal(writeUpInfos)
		if err != nil {
			log.Panicln("writeup Marshal error: ", err)
			return
		}
		_, err = w.Write(infosJson)
		if err != nil {
			log.Println("writeup infos send error: ", err)
			return
		}

	}

}
