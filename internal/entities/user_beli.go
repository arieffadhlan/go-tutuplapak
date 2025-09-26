package entities

type UserBeli struct {
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
	IsAdmin  bool   `db:"isAdmin"`
}
