package user

import (
	"mime/multipart"
	"net/http"
)

func WirteUpHandle(w http.ResponseWriter, req *http.Request, userId uint64) {
	if req.Method != "POST" {
		w.WriteHeader(404)
		w.Write(nil)
		return
	}
	err := req.ParseMultipartForm(16 << 20)
	if err != nil {
		return
	} // TODO: move this to config
	var file multipart.File
	var fileHandle *multipart.FileHeader
	file, fileHandle, err = req.FormFile("writeup")
	if err != nil {
		return
	}
	defer file.Close()
	// TODO: the real upload logic
}
