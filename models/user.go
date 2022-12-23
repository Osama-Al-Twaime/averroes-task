package models

type User struct {
	Id       int    `json:"id" sql:"id"`
	Name     string `json:"name" sql:"name"`
	Email    string `json:"email" sql:"email"`
	Age      int    `json:"age" sql:"age"`
	Password string `sql:"password"`
}

type UserRegisterRequest struct {
	Name     string `json:"name" sql:"name"`
	Email    string `json:"email" sql:"email"`
	Age      int    `json:"age" sql:"age"`
	Password string `sql:"password"`
}
