package model

type VideoListItem struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Duration int `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

type VideoItem struct {
	Item VideoListItem
	Url  string `json:"url"`
}
