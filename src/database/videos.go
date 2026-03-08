package database

import (
	"database/sql"
	"time"
)

type Video struct {
	ID        int
	OwnerID   int
	OwnerName string
	Title     string
	VideoHash string
	Extension string
	IsPrivate bool
	CreatedAt time.Time
	IsOwner   bool
}

func InsertVideo(ownerID int, title, videoHash string, extension string, isPrivate bool) (int, error) {
	var id int
	err := database.QueryRow(`
		INSERT INTO videos (owner_id, title, video_hash, extension, is_private)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, ownerID, title, videoHash, extension, isPrivate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetVideoByID(videoID int, currentUserID int) (*Video, error) {
	var video Video

	err := database.QueryRow(`
        SELECT v.id, v.owner_id, u.username, v.title, v.video_hash, v.extension, v.is_private, v.created_at, 
               (v.owner_id = $2) AS is_owner
        FROM videos v
        JOIN users u ON v.owner_id = u.id
        WHERE v.id = $1
    `, videoID, currentUserID).Scan(
		&video.ID,
		&video.OwnerID,
		&video.OwnerName,
		&video.Title,
		&video.VideoHash,
		&video.Extension,
		&video.IsPrivate,
		&video.CreatedAt,
		&video.IsOwner,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &video, nil
}
