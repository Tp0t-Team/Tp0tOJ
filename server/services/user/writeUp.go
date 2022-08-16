package user

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"server/services/database/resolvers"
	"server/utils/configure"
	"strconv"
	"strings"
	"time"
)

func WriteUpHandle(w http.ResponseWriter, req *http.Request, userId uint64) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
		return
	}
	err := req.ParseMultipartForm(16 << 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	} // TODO: move this to config
	var file multipart.File
	var fileHandle *multipart.FileHeader
	file, fileHandle, err = req.FormFile("writeup")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}

	defer file.Close()
	fileNameParts := strings.Split(fileHandle.Filename, ".")
	extname := fileHandle.Filename[len(fileNameParts[0]):]
	// TODO: add more check, remove the old one
	outFile, err := os.Create(configure.WriteUpPath + "/" + strconv.FormatUint(userId, 10) + extname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	resolvers.BehaviorUploadWP(userId, time.Now(), nil)
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
	return
}
