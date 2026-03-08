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
