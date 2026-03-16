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

	accessToken, newRefreshToken, err := rotateTokens(userID, refreshToken)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(ctx, newRefreshToken)

	return ctx.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func parseRefreshRequest(ctx fiber.Ctx) (string, error) {
	refreshToken := getRefreshTokenCookie(ctx)
	if refreshToken == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "refresh token cookie is required")
	}

	return refreshToken, nil
}

func rotateTokens(userID int, oldRefreshToken string) (string, string, error) {
	accessToken, err := GenerateToken(userID)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't generate token for %v", userID))
		return "", "", fiber.ErrInternalServerError
	}

	newRefreshToken, err := database.RotateRefreshToken(userID, oldRefreshToken, 7*24*time.Hour)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't rotate refreshtoken for %v", userID))
		return "", "", fiber.ErrInternalServerError
	}

	return accessToken, newRefreshToken, nil
}
