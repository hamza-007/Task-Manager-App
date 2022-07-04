package services

import (
	"database/sql"
	"errors"

	"github.com/hamza-007/Task-Manager-App/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	AddUser(*models.User) (error)
	GetUser(string) (*models.User,error)
	UpdateUser(string) (models.User,error)
}

type UserSvc struct{
	BD 	*sql.DB
}

func NewUserService(bd *sql.DB) UserService {
	return &UserSvc{
		BD: bd,
	}
}
func (us *UserSvc)AddUser(usr *models.User) error{
	var email string	
	password, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14)
	if err != nil { return err }
	row, err := us.BD.Query("SELECT email FROM users WHERE email = ?", &usr.Email)
	if err != nil {
		return  err
	}
	for row.Next() {
		err := row.Scan(&email)
		if err != nil {
			return err
		}
	}

	if email == usr.Email {
		return errors.New("user already exists !!!")
	}
	stmt, err := us.BD.Prepare("INSERT INTO users(userid,username,email,passwrd) VALUES (?,?,?,?) ")
	if err != nil {
		return err 
	}
	_, err = stmt.Exec(&usr.Id, &usr.Username, &usr.Email,password)
	return err
}

func (us *UserSvc)GetUser(id string) (*models.User,error){		
	var user models.User	
	
	row, err := us.BD.Query("SELECT * FROM users WHERE email = ?", &id)
	if err != nil {
		return  nil,err
	}
	for row.Next() {
		err := row.Scan(&user.Id,&user.Username,&user.Email,&user.Password)
		if err != nil {
			return nil,err
		}
	}

	if user.Email == "" {
		return nil,errors.New("user not found !!!")
	}
	return &user,nil

}

func (us *UserSvc)UpdateUser(id string) (models.User,error){
	var user models.User
	return user,nil 
}