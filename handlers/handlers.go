package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "mime/multipart"
	_ "io"
	_ "os"
	"github.com/AlK2x/simple_video_service/packages/repository"
)


func Router(repo *repository.VideoRepository) http.Handler {
	l := ListHandler{Db: repo}
	v := VideoHandler{Db: repo}
	u := UploadHandler{Db: repo}
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/list", l.list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", v.video).Methods(http.MethodGet)
	s.HandleFunc("/video", u.upload).Methods(http.MethodPost)

	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"url": r.URL,
			"removeAddr": r.RemoteAddr,
			"userAgent": r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello Worldj")
}

