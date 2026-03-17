package database

import (
	"fmt"
	"simon-weij/wayland-recorder-backend/src/dto"

	"gorm.io/gorm"
)

func InsertUserIntoDatabase(user dto.User) (int, error) {
	newUser := UserModel{
		Username:     user.Username,
		PasswordHash: user.Password,
	}

	err := database.Create(&newUser).Error
	if err != nil {
		return 0, err
	}

	return newUser.ID, nil
}

func ValueAlreadyExists(whatExists string, value string) (bool, error) {
	if whatExists != "username" {
		return false, fmt.Errorf("unsupported field: %s", whatExists)
	}

	var count int64
	err := database.Model(&UserModel{}).Where("username = ?", value).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserAuthByUsername(username string) (*dto.UserAuth, error) {
	var user dto.UserAuth

	type userAuthQueryResult struct {
		ID           int    `gorm:"column:id"`
		PasswordHash string `gorm:"column:password_hash"`
	}

	var result userAuthQueryResult
	err := database.Model(&UserModel{}).
		Select("id", "password_hash").
		Where("username = ?", username).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	user.ID = result.ID
	user.PasswordHash = result.PasswordHash

	return &user, nil
}
