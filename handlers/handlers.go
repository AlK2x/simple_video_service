package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "mime/multipart"
	_ "io"
	_ "os"
	"io"
	"os"
)

type VideoListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type VideoInfo struct {
	Item VideoListItem
	Url  string `json:"url"`
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/list", list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", video).Methods(http.MethodGet)
	s.HandleFunc("/video", uploadVideo).Methods(http.MethodPost)

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

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	inputFile, header, err := r.FormFile("file[]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType != "video/mp4" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fileName := header.Filename

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	bytes, err := io.Copy(file, inputFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, bytes)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello Worldj")
}

