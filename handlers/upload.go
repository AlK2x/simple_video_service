package handlers

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"github.com/satori/go.uuid"
	"github.com/AlK2x/simple_video_service/packages/repository"
	"github.com/AlK2x/simple_video_service/packages/model"
)

type UploadHandler struct {
	Db *repository.VideoRepository
}

func (u UploadHandler) upload(w http.ResponseWriter, r *http.Request) {
	inputFile, header, err := r.FormFile("file[]")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType != "video/mp4" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fileName := header.Filename

	u1, err := uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := "content/" + u1.String() + "/"
	fileUrl := path + fileName
	os.MkdirAll(path, os.ModePerm)
	file, err := os.OpenFile(fileUrl, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	bytes, err := io.Copy(file, inputFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	video := &model.VideoItem{
		Item: model.VideoListItem{
			Id: u1.String(),
			Name: fileName,
			Duration: 42,
			Thumbnail: "",
		},
		Url: fileUrl,
	}
	err = u.Db.SaveVideo(video)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, bytes)
}

