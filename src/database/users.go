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

package database

import (
	"fmt"
	"simon-weij/wayland-recorder-backend/src/dto"
)

func InsertUserIntoDatabase(user dto.User) (int, error) {
	var lastInsertId int
	query := `INSERT INTO users (username, password_hash) 
              VALUES ($1, $2) 
              RETURNING id`

	err := database.QueryRow(query, user.Username, user.Password).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func ValueAlreadyExists(whatExists string, value string) (bool, error) {
	var valueExists bool

	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE %s = $1
		)
	`, whatExists)
	err := database.QueryRow(query, value).Scan(&valueExists)
	return valueExists, err
}

func GetUserAuthByUsername(username string) (*dto.UserAuth, error) {
	var user dto.UserAuth

	err := database.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
