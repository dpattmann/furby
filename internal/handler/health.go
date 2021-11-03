package handler

import (
	"net/http"
)

func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Encoding", "UTF-8")

		w.Write([]byte("Ok"))
	})
}
