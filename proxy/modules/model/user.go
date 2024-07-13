package model

type User struct {
	Email    string `json:"email" binding:"required" example:"admin@example.com"`
	Password string `json:"password" binding:"required" example:"password"`
}
