package auth

import (
	"fmt"
	"simon-weij/wayland-recorder-backend/src/database"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// /auth/refresh
func RefreshToken(ctx fiber.Ctx) error {
	refreshToken, err := parseRefreshRequest(ctx)
	if err != nil {
		return err
	}

	userID, err := database.GetUserIDFromRefreshToken(refreshToken)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	accessToken, newRefreshToken, err := rotateTokens(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func parseRefreshRequest(ctx fiber.Ctx) (string, error) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.Bind().Body(&body); err != nil {
		return "", fiber.ErrBadRequest
	}

	if body.RefreshToken == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "refresh_token is required")
	}

	return body.RefreshToken, nil
}

func rotateTokens(userID int) (string, string, error) {
	accessToken, err := GenerateToken(userID)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't generate token for %v", userID))
		return "", "", fiber.ErrInternalServerError
	}

	newRefreshToken, err := database.CreateRefreshToken(userID, 7*24*time.Hour)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't generate refreshtoken for %v", userID))
		return "", "", fiber.ErrInternalServerError
	}

	return accessToken, newRefreshToken, nil
}
