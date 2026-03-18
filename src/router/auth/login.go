package auth

import (
	"errors"
	"fmt"
	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/dto"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// /auth/login
func Login(ctx fiber.Ctx) error {
	body, err := parseAndValidate(ctx)
	if err != nil {
		return err
	}

	user, err := authenticate(body)
	if err != nil {
		return err
	}

	token, refreshToken, err := issueTokens(user.ID)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(ctx, refreshToken)

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func parseAndValidate(ctx fiber.Ctx) (*LoginRequest, error) {
	var body LoginRequest

	if err := ctx.Bind().Body(&body); err != nil {
		return nil, fiber.ErrBadRequest
	}

	if body.Username == "" || body.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Username and password are required")
	}

	return &body, nil
}

func authenticate(body *LoginRequest) (*dto.UserAuth, error) {
	user, err := database.GetUserAuthByUsername(body.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
		}
		log.Warn(fmt.Sprintf("Couldn't get the user of %s", body.Username))
		return nil, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(body.Password),
	); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	}

	return user, nil
}
func issueTokens(userID int) (string, string, error) {
	token, err := GenerateToken(userID)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't generate token for %v", userID))
		return "", "", fiber.ErrInternalServerError
	}

	refresh_token, err := database.CreateRefreshToken(userID, 7*24*time.Hour)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't create refresh token for %v with error %v", userID, err))
		return "", "", fiber.ErrInternalServerError
	}

	return token, refresh_token, nil
}
