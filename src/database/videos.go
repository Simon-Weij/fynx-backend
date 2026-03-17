package database

import (
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
	record := VideoModel{
		OwnerID:   ownerID,
		Title:     title,
		VideoHash: videoHash,
		Extension: extension,
		IsPrivate: isPrivate,
	}

	err := database.Create(&record).Error
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

func GetVideoByID(videoID int, currentUserID int) (*Video, error) {
	var video Video
	result := database.Table("videos v").
		Select(`v.id, v.owner_id, u.username as owner_name, v.title, v.video_hash, v.extension, v.is_private, v.created_at, (v.owner_id = ?) as is_owner`, currentUserID).
		Joins("JOIN users u ON v.owner_id = u.id").
		Where("v.id = ?", videoID).
		Scan(&video)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &video, nil
}
