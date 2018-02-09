package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

type VideListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type PlayVideoInfo struct {
	Item VideListItem
	Url  string `json:"url"`
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/list", list).Methods(http.MethodGet)
	s.HandleFunc("/video/{ID}", video).Methods(http.MethodGet)

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

