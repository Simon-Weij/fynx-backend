package user

import (
	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/router/auth"

	"github.com/gofiber/fiber/v3"
)

// /api/user/data
func UserData(ctx fiber.Ctx) error {
	userID, err := auth.GetUID(ctx)
	if err != nil {
		return err
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{
		"username": user.Username,
	})
}
