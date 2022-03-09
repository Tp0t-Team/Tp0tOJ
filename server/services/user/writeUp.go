package user

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func WriteUpHandle(w http.ResponseWriter, req *http.Request, userId uint64) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		w.Write(nil)
		return
	}
	//TODO: need login and role check
	err := req.ParseMultipartForm(16 << 20)
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		log.Println(err)
		return
	} // TODO: move this to config
	var file multipart.File
	var fileHandle *multipart.FileHeader
	file, fileHandle, err = req.FormFile("writeup")
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		log.Println(err)
		return
	}
	defer file.Close()
	fileNameParts := strings.Split(fileHandle.Filename, ".")
	extname := fileNameParts[len(fileNameParts)-1]
	// TODO: add more check
	outFile, err := os.Open(strconv.FormatUint(userId, 10) + extname)
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		log.Println(err)
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
		log.Println(err)
		return
	}
	w.WriteHeader(200)
	w.Write(nil)
	return
}
