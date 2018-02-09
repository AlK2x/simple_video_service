package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/gorilla/mux"
	"fmt"
)

func video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	responseVideo := VideoInfo{
		Item: VideoListItem{
			Id:        "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Name:      "Black Retrospective Woman",
			Duration:  15,
			Thumbnail: fmt.Sprintf("/content/%s/screen.jpg", id),
		},
		Url: fmt.Sprintf("/content/%s/index.mp4", id),
	}
	data, err := json.Marshal(responseVideo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(w, string(data))
}
