package handlers

import (
	"net/http"
	"encoding/json"
	"io"
	"github.com/gorilla/mux"
	"database/sql"
)

type VideoHandler struct {
	Db *sql.DB
}

func (v VideoHandler) video(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]

	responseVideo := VideoItem{}
	stmt, err := v.Db.Prepare(`SELECT video_key, title, duration, thumbnail_url, url FROM video WHERE video_key = ?`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(
		&responseVideo.Item.Id,
		&responseVideo.Item.Name,
		&responseVideo.Item.Duration,
		&responseVideo.Item.Thumbnail,
		&responseVideo.Url,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

/*	responseVideo := VideoItem{
		Item: VideoListItem{
			Id:        "d290f1ee-6c54-4b01-90e6-d701748f0851",
			Name:      "Black Retrospective Woman",
			Duration:  15,
			Thumbnail: fmt.Sprintf("/content/%s/screen.jpg", id),
		},
		Url: fmt.Sprintf("/content/%s/index.mp4", id),
	}*/
	data, err := json.Marshal(responseVideo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	io.WriteString(w, string(data))
}
