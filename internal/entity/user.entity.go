package entity

type User struct {
	UserID   string `db:"user_id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
}
