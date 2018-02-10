package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"database/sql"
)

type ListHandler struct {
	Db *sql.DB
}

func (l ListHandler) list(w http.ResponseWriter, _ *http.Request) {

	rows, err := l.Db.Query(`SELECT video_key, title, duration, thumbnail_url FROM video`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	videos := []VideoListItem{}
	for rows.Next() {
		var videoItem VideoListItem
		err :=  rows.Scan(&videoItem.Id, &videoItem.Name, &videoItem.Duration, &videoItem.Thumbnail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, videoItem)
	}

	r, err := json.Marshal(videos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json: charset=UTF-8")
	io.WriteString(w, string(r))
}