package main

import (
	"net/http"
	"github.com/AlK2x/simple-video-service/handlers"
	log "github.com/Sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the serve")
	router := handlers.Router()
	log.Fatal(http.ListenAndServe(serverUrl, router))
}
