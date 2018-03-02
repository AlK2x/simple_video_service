package repository

import (
	"database/sql"
	"github.com/AlK2x/simple_video_service/packages/model"
	"net/http"
)

type VideoRepository struct {
	db *sql.DB
}

func CreateVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db}
}

func (v VideoRepository) GetVideos() ([]model.VideoListItem, error) {
	rows, err := v.db.Query(`SELECT video_key, title, duration, thumbnail_url FROM video`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := []model.VideoListItem{}
	for rows.Next() {
		var videoItem model.VideoListItem
		err :=  rows.Scan(&videoItem.Id, &videoItem.Name, &videoItem.Duration, &videoItem.Thumbnail)
		if err != nil {
			return nil, err
		}
		videos = append(videos, videoItem)
	}
	return videos, nil
}

func (v VideoRepository) GetVideo(id string) (*model.VideoItem, error) {
	responseVideo := model.VideoItem{}
	stmt, err := v.db.Prepare(`SELECT video_key, title, duration, thumbnail_url, url FROM video WHERE video_key = ?`)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &responseVideo, nil
}

func (v VideoRepository) SaveVideo(video *model.VideoItem) error {
	q := `INSERT INTO video (video_key, title, status, duration, url) VALUES(?, ?, ?, ?, ?)`
	stmt, err := v.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(video.Item.Id, video.Item.Name, 1, video.Item.Duration, video.Url)
	if err != nil {
		return err
	}
	return nil
}

func (v VideoRepository) UpdateVideo(key string, duration float64, thumbnail string) error {
	q := `UPDATE video SET duration = ?, thumbnail_url = ? WHERE video_key = ?`
	stmt, err := v.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(duration, thumbnail)
	if err != nil {
		return err
	}
	return nil
}

func (v VideoRepository) GetUnprocessedVideo() (*model.VideoItem, error) {
	responseVideo := model.VideoItem{}
	stmt, err := v.db.Prepare(`SELECT video_key, title, duration, thumbnail_url, url FROM video WHERE status = 1 LIMIT 1`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow().Scan(
		&responseVideo.Item.Id,
		&responseVideo.Item.Name,
		&responseVideo.Item.Duration,
		&responseVideo.Item.Thumbnail,
		&responseVideo.Url,
	)

	if err != nil {
		return nil, err
	}

	return &responseVideo, nil
}
