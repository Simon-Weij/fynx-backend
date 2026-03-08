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
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type DatabaseCredentials struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
}

var database *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func InitialiseDatabase() {
	godotenv.Load()
	credentials := DatabaseCredentials{
		dbUser:     os.Getenv("DB_USERNAME"),
		dbPassword: os.Getenv("DB_PASSWORD"),
		dbHost:     os.Getenv("DB_HOST"),
		dbPort:     os.Getenv("DB_PORT"),
		dbName:     os.Getenv("DB_NAME"),
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		credentials.dbUser,
		credentials.dbPassword,
		credentials.dbHost,
		credentials.dbPort,
		credentials.dbName,
	)

	var err error
	database, err = sql.Open("pgx", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	setupUsersTable()
	setupRefreshTokenTable()
	setupVideosTable()
}

func setupUsersTable() {
	database.Query(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`)
}

func setupVideosTable() {
	database.Query(`
	CREATE TABLE IF NOT EXISTS videos (
		id SERIAL PRIMARY KEY,
		owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		video_hash VARCHAR(255) NOT NULL,
		extension VARCHAR(10) NOT NULL,
		is_private BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`)
}

func setupRefreshTokenTable() {
	database.Query(`
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		hashed_token VARCHAR(255) NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`)
}
