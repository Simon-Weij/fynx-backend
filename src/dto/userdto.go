package dto

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password_hash"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAuth struct {
	ID           int
	PasswordHash string
}
