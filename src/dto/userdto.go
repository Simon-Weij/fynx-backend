package dto

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password_hash"`
}

type UserAuth struct {
	ID           int
	PasswordHash string
}
