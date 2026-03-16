package auth

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

const refreshTokenCookieName = "refresh_token"

func setRefreshTokenCookie(ctx fiber.Ctx, token string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   isCookieSecure(),
		SameSite: "lax",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func getRefreshTokenCookie(ctx fiber.Ctx) string {
	return ctx.Cookies(refreshTokenCookieName)
}

func clearRefreshTokenCookie(ctx fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   isCookieSecure(),
		SameSite: "lax",
		Expires:  time.Unix(0, 0),
	})
}

func isCookieSecure() bool {
	value := os.Getenv("COOKIE_SECURE")
	return value == "1"
}
