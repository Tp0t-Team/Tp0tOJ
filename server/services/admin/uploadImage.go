package admin

import (
	"log"
	"mime/multipart"
	"net/http"
	"server/services/kube"
	"strings"
)

func UploadImage(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
		return
	}
	err := req.ParseMultipartForm(16 << 20) // TODO: set bigger?
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	var file multipart.File
	//var fileHandle *multipart.FileHeader
	name := req.FormValue("name")
	file, _, err = req.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	defer file.Close()

	err = kube.ImgBuild(file, strings.ToLower(name))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
