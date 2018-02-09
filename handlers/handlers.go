package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"io"
)

type VideoInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type PlayVideoInfo struct {
	Item VideoInfo
	Url string `json:"url"`
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/list", getVideoPreview).Methods(http.MethodGet)
	s.HandleFunc("/video/d290f1ee-6c54-4b01-90e6-d701748f0851", getVideoContent).Methods(http.MethodGet)

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

func getVideoPreview(w http.ResponseWriter, _ *http.Request) {
	response := []VideoInfo{{
		Id: "d290f1ee-6c54-4b01-90e6-d701748f0851",
		Name: "Black Retrospective Woman",
		Duration: 15,
		Thumbnail: "/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
	}}

	r, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json: charset=UTF-8")
	io.WriteString(w, string(r))
}

func getVideoContent(w http.ResponseWriter, _ *http.Request) {

	responseVideo := PlayVideoInfo{
		Item: VideoInfo{
			Id:        "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Name:      "Black Retrospective Woman",
			Duration:  15,
			Thumbnail: "/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
		},
		Url: "/content/d290f1ee-6c54-4b01-90e6-d701748f0851/index.mp4",
	}
	r, err := json.Marshal(responseVideo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(w, string(r))
}