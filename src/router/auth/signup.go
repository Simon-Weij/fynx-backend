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

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
)

// /auth/signup
func Signup(ctx fiber.Ctx) error {
	userDto, err := parseAndValidateSignup(ctx)
	if err != nil {
		return err
	}

	if err := checkUserConflicts(userDto); err != nil {
		return err
	}

	if err := processAndStoreUser(userDto); err != nil {
		return err
	}

	return nil
}

func parseAndValidateSignup(ctx fiber.Ctx) (*dto.User, error) {
	var userDto dto.User

	if err := ctx.Bind().Body(&userDto); err != nil {
		return nil, fiber.ErrBadRequest
	}

	return &userDto, nil
}

func checkUserConflicts(userDto *dto.User) error {
	if err := checkForConflicts("username", userDto.Username); err != nil {
		return err
	}
	return nil
}

func processAndStoreUser(userDto *dto.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Warn("Couldn't hash the password!")
		return fiber.NewError(fiber.ErrInternalServerError.Code, "couldn't hash password")
	}
	userDto.Password = string(hash)

	database.InsertUserIntoDatabase(*userDto)
	return nil
}

func checkForConflicts(whatExists string, value string) error {
	userWithEmailExists, err := database.ValueAlreadyExists(whatExists, value)
	if err != nil {
		log.Warn(fmt.Sprintf("Couldn't check if %s already exists", whatExists))
		return fiber.NewError(fiber.ErrInternalServerError.Code, "couldn't check if %s already exists", whatExists)
	}
	if userWithEmailExists {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("user with %s %s already exists.", whatExists, value))
	}

	return nil
}
