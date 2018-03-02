package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/gorilla/mux"
	"github.com/AlK2x/simple_video_service/packages/repository"
)

type VideoHandler struct {
	Db *repository.VideoRepository
}

func (v VideoHandler) video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	responseVideo, err := v.Db.GetVideo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(responseVideo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(w, string(data))
}
