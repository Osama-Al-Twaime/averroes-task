package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Age      int    `json:"age" db:"age"`
	Password string `db:"password"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Age      int    `json:"age" db:"age"`
	Password string `json:"password" db:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
