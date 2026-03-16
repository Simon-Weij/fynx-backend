package database

import (
	"database/sql"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	query := `
		INSERT INTO refresh_tokens (user_id, hashed_token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var tokenID int
	err = database.QueryRow(query, userID, hashedToken, expiresAt).Scan(&tokenID)
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return rawToken, nil
}

func GetUserIDFromRefreshToken(rawToken string) (int, error) {
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	var userID int
	var expiresAt time.Time

	err := database.QueryRow(
		`SELECT user_id, expires_at FROM refresh_tokens WHERE hashed_token = $1`,
		hashedToken,
	).Scan(&userID, &expiresAt)

	if err != nil {
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("refresh token expired")
	}

	return userID, nil
}

func RefreshToken(userID int, rawRefreshToken string, tokenDuration time.Duration) (string, error) {
	hash := sha256.Sum256([]byte(rawRefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	var expiresAt time.Time
	query := `
		SELECT expires_at
		FROM refresh_tokens
		WHERE user_id = $1 AND hashed_token = $2
	`
	err := database.QueryRow(query, userID, hashedToken).Scan(&expiresAt)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(expiresAt) {
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

	_, err := database.Exec(
		`DELETE FROM refresh_tokens WHERE hashed_token = $1`,
		hashedToken,
	)

	return err
}

func RotateRefreshToken(userID int, oldRawToken string, duration time.Duration) (string, error) {
	oldHashedToken := hashRefreshToken(oldRawToken)
	newRawToken, err := generateRawRefreshToken()
	if err != nil {
		return "", err
	}

	newHashedToken := hashRefreshToken(newRawToken)
	newExpiresAt := time.Now().Add(duration)

	tx, err := database.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := deleteRefreshTokenInTx(tx, userID, oldHashedToken)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	err = insertRefreshTokenInTx(tx, userID, newHashedToken, newExpiresAt)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
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

func deleteRefreshTokenInTx(tx *sql.Tx, userID int, hashedToken string) (sql.Result, error) {
	return tx.Exec(
		`DELETE FROM refresh_tokens WHERE user_id = $1 AND hashed_token = $2`,
		userID,
		hashedToken,
	)
}

func insertRefreshTokenInTx(tx *sql.Tx, userID int, hashedToken string, expiresAt time.Time) error {
	_, err := tx.Exec(
		`INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3)`,
		userID,
		hashedToken,
		expiresAt,
	)

	return err
}
