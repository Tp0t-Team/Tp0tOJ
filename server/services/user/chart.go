package user

import "net/http"

func ChartHandle(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(nil)
		return
	}

}
