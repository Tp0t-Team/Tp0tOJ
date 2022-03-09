package admin

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/services/database/resolvers"
	"strconv"
	"strings"
)

const writeUpPath string = "./writeup"

func DownloadAllWP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		w.Write(nil)
		return
	}

	//TODO: pack writeup folder and download
	var zipFileName = "WP.zip"
	var writeUpPath string
	err := os.RemoveAll(zipFileName)
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
		return
	}
	zipFile, _ := os.Create(zipFileName)
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Panicln("DownloadAllWP: ", err)
			return
		}
	}(zipFile)

	archive := zip.NewWriter(zipFile)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			log.Panicln("DownloadAllWP: ", err)
			return
		}
	}(archive)
	err = filepath.Walk(writeUpPath, func(path string, info os.FileInfo, _ error) error {
		if path == writeUpPath {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, writeUpPath+`\`)
		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Panicln("DownloadAllWP: ", err)
					return
				}
			}(file)
			_, err := io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
		return
	}
	//TODO: send zip

}
func DownloadWPByUserId() {

}

type WriteUpInfoResult struct {
	UserId string
	Name   string
	Mail   string
	Solved int //solved challenge number
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
			UserId: string(userId),
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
