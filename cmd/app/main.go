package main

import (
	"net/http"
	"github.com/AlK2x/simple_video_service/handlers"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/AlK2x/simple_video_service/packages/repository"

)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	db, err := sql.Open("mysql", `root:12345@/simple_video_service`)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	killSignalChan := getKillSignalChan()
	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the serve")

	repository := repository.CreateVideoRepository(db)
	srv := startServer(serverUrl, repository)
	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}

func startServer(serverUrl string, db *repository.VideoRepository) *http.Server {
	router := handlers.Router(db)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT")
	case syscall.SIGTERM:
		log.Info("got SIGTERM")
	}
}
