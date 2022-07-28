package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hamza-007/Task-Manager-App/handlers"
	"github.com/hamza-007/Task-Manager-App/models"
	"github.com/hamza-007/Task-Manager-App/services"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userservice services.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	user.Id = uuid.New().String()

	err = uc.UserService.VerifUser(&user.Email)
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	err = uc.UserService.AddUser(&user)

	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var res = handlers.NewHTTPResponse(http.StatusCreated, "user saved succesfully !!")
	c.JSON(http.StatusCreated, res)
	return
}

func (uc *UserController) Login(c *gin.Context) {

	var payload models.UserLogin
	err := c.Bind(&payload)
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, "errroooor")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	user, err := uc.UserService.GetUserByEmail(payload.Email)
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusNotFound, err)
		c.JSON(http.StatusNotFound, res)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.SetCookie("jwt", token, 1000*60*60*24, "/", "", false, false)

	var res = handlers.NewHTTPResponse(http.StatusAccepted, "logged in succesfully ! ")
	c.JSON(http.StatusAccepted, res)
	return
}

func (uc *UserController) User(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user, err := uc.UserService.GetUserById(claims.Issuer)
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	var res = handlers.NewHTTPResponse(http.StatusCreated, user)
	c.JSON(http.StatusCreated, res)
	return
}

func (uc *UserController) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1000, "/", "", false, false)

	var res = handlers.NewHTTPResponse(http.StatusOK, "logout succes")
	c.JSON(http.StatusOK, res)
	return
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userrouter := rg.Group("/user")
	userrouter.POST("/create", uc.Register)
	userrouter.POST("/login", uc.Login)
	userrouter.GET("/", uc.User)
	userrouter.GET("/logout", uc.Logout)
	// taskrouter.DELETE("/delete/:id",tc.DeleteTask)
}
