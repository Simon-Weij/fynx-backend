package auth

import (
	"simon-weij/wayland-recorder-backend/src/database"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// /auth/logout
func Logout(ctx fiber.Ctx) error {
	refreshToken := getRefreshTokenCookie(ctx)
	if refreshToken != "" {
		if err := database.DeleteRefreshToken(refreshToken); err != nil {
			log.Warnf("Couldn't delete refresh token: %v", err)
			return fiber.ErrInternalServerError
		}
	}

	clearRefreshTokenCookie(ctx)

	return ctx.SendStatus(fiber.StatusNoContent)
}
