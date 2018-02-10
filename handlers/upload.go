package handlers

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"database/sql"
	"github.com/satori/go.uuid"
)

type UploadHandler struct {
	Db *sql.DB
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

	q := `INSERT INTO video (video_key, title, status, duration, url) VALUES(?, ?, ?, ?, ?)`
	stmt, err := u.Db.Prepare(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u1.String(), "New video", 3, 42, fileUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, bytes)
}

