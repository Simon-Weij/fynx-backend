package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateRefreshToken(userID int, duration time.Duration) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	rawToken := hex.EncodeToString(b)

	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(duration)

	refreshToken := RefreshTokenModel{
		UserID:      userID,
		HashedToken: hashedToken,
		ExpiresAt:   expiresAt,
	}

	err = database.Create(&refreshToken).Error
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return rawToken, nil
}

func GetUserIDFromRefreshToken(rawToken string) (int, error) {
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	var tokenRecord RefreshTokenModel
	err := database.Where("hashed_token = ?", hashedToken).First(&tokenRecord).Error

	if err != nil {
		return 0, err
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return 0, errors.New("refresh token expired")
	}

	return tokenRecord.UserID, nil
}

func RefreshToken(userID int, rawRefreshToken string, tokenDuration time.Duration) (string, error) {
	hash := sha256.Sum256([]byte(rawRefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	var tokenRecord RefreshTokenModel
	err := database.Where("user_id = ? AND hashed_token = ?", userID, hashedToken).First(&tokenRecord).Error
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return "", errors.New("refresh token expired")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(tokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func DeleteRefreshToken(rawToken string) error {
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	return database.Where("hashed_token = ?", hashedToken).Delete(&RefreshTokenModel{}).Error
}

func RotateRefreshToken(userID int, oldRawToken string, duration time.Duration) (string, error) {
	oldHashedToken := hashRefreshToken(oldRawToken)
	newRawToken, err := generateRawRefreshToken()
	if err != nil {
		return "", err
	}

	newHashedToken := hashRefreshToken(newRawToken)
	newExpiresAt := time.Now().Add(duration)

	err = database.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND hashed_token = ?", userID, oldHashedToken).Delete(&RefreshTokenModel{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Create(&RefreshTokenModel{
			UserID:      userID,
			HashedToken: newHashedToken,
			ExpiresAt:   newExpiresAt,
		}).Error
	})
	if err != nil {
		return "", err
	}

	return newRawToken, nil
}

func hashRefreshToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

func generateRawRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	return hex.EncodeToString(b), nil
}
