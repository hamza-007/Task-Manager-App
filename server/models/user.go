package models

type User struct {
	Id       	string 		`json:"id"`
	Username 	string 		`json:"username" binding:"required"`
	Email   	string   	`json:"email" binding:"required"`
	Password 	string 		`json:"password" binding:"required"`
}

type UserLogin struct {
	Email   	string   	`json:"email" binding:"required"`
	Password 	string 		`json:"password" binding:"required"`
}