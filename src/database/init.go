package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseCredentials struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbName     string
}

var database *gorm.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserModel struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"size:50;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

type VideoModel struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	OwnerID   int       `gorm:"column:owner_id;not null;index"`
	Title     string    `gorm:"size:255;not null"`
	VideoHash string    `gorm:"column:video_hash;size:255;not null"`
	Extension string    `gorm:"size:10;not null"`
	IsPrivate bool      `gorm:"column:is_private;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (VideoModel) TableName() string {
	return "videos"
}

type RefreshTokenModel struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	UserID      int       `gorm:"column:user_id;not null;index"`
	HashedToken string    `gorm:"column:hashed_token;size:255;not null;uniqueIndex"`
	ExpiresAt   time.Time `gorm:"column:expires_at;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RefreshTokenModel) TableName() string {
	return "refresh_tokens"
}

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
	database, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	err = database.AutoMigrate(&UserModel{}, &RefreshTokenModel{}, &VideoModel{})
	if err != nil {
		log.Fatal(err)
	}
}
