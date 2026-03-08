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
