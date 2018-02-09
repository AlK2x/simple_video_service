package handlers

import (
	"net/http"
	"encoding/json"
	"io"
)

func list(w http.ResponseWriter, _ *http.Request) {
	response := []VideoListItem{{
		Id: "d290f1ee-6c54-4b01-90e6-d701748f0851",
		Name: "Black Retrospective Woman",
		Duration: 15,
		Thumbnail: "/content/d290f1ee-6c54-4b01-90e6-d701748f0851/screen.jpg",
	},
	{
		Id: "hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345",
		Name: "Go Really TEASER-HD",
		Duration: 41,
		Thumbnail: "/content/hjkhhjk3-23j4-j45k-erkj-kj3k4jl2k345/screen.jpg",
	},
	{
		Id: "sldjfl34-dfgj-523k-jk34-5jk3j45klj34",
		Name: "Танцор",
		Duration: 92,
		Thumbnail: "/content/sldjfl34-dfgj-523k-jk34-5jk3j45klj34/screen.jpg",
	}}

	r, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json: charset=UTF-8")
	io.WriteString(w, string(r))
}