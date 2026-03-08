/*
 * Wayland recorder is a way to easily make clips and share them.
 * Copyright (C) 2026 Simon-Weij
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package auth

import (
	"fmt"
	"simon-weij/wayland-recorder-backend/src/database"
	"simon-weij/wayland-recorder-backend/src/dto"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
)

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

	return ctx.JSON(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func parseAndValidate(ctx fiber.Ctx) (*dto.LoginRequest, error) {
	var body dto.LoginRequest

	if err := ctx.Bind().Body(&body); err != nil {
		return nil, fiber.ErrBadRequest
	}

	if body.Username == "" || body.Password == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Username and password are required")
	}

	return &body, nil
}

func authenticate(body *dto.LoginRequest) (*dto.UserAuth, error) {
	user, err := database.GetUserAuthByUsername(body.Username)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't get the user of %s", body.Username))
		return nil, fiber.ErrInternalServerError
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(body.Password),
	); err != nil {
		log.Warn("Couldn't compare hash")
		return nil, fiber.ErrInternalServerError
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
