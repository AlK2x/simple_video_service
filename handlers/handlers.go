package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "mime/multipart"
	_ "io"
	_ "os"
	"database/sql"
)

type VideoListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type VideoItem struct {
	Item VideoListItem
	Url  string `json:"url"`
}

func Router(db *sql.DB) http.Handler {
	l := ListHandler{Db: db}
	v := VideoHandler{Db: db}
	u := UploadHandler{Db: db}
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

