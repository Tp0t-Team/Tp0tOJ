package user

import (
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	// remove the old writeup
	fileName := strconv.FormatUint(userId, 10)
	removeList := []string{}
	err = filepath.Walk(configure.WriteUpPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.Split(info.Name(), ".")[0] == fileName {
			removeList = append(removeList, path)
		}
		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	for _, f := range removeList {
		err := os.Remove(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			log.Println(err)
			return
		}
	}
	outFile, err := os.Create(configure.WriteUpPath + "/" + fileName + extname)
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
