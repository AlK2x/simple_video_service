package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/AlK2x/simple_video_service/packages/repository"
)

type ListHandler struct {
	Db *repository.VideoRepository
}

func (l ListHandler) list(w http.ResponseWriter, _ *http.Request) {
	videos, err := l.Db.GetVideos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r, err := json.Marshal(videos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json: charset=UTF-8")
	io.WriteString(w, string(r))
}