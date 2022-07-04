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

type UserController struct{
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


	err = uc.UserService.AddUser(&user)

	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	var res = handlers.NewHTTPResponse(http.StatusCreated,"user saved succesfully !!")
	c.JSON(http.StatusCreated,res)
	return
}

func (uc *UserController) Login(c *gin.Context)  {
	
	var payload models.UserLogin
	err := c.Bind(&payload); 
	if err != nil {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, "errroooor")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	user,err := uc.UserService.GetUser(payload.Email) 
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusNotFound,err)
		c.JSON(http.StatusNotFound,res)
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
	c.SetCookie("jwt",token, 10*1000, "/", "", false, true)
	
}
/*
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("info", fiber.Map{
			"Message": "not authentificated ",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User1

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.Render("home",user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Render("info", fiber.Map{
		"Message": "logout succes  ",
	})
}*/


func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	taskrouter := rg.Group("/user")
	taskrouter.POST("/create", uc.Register)
	taskrouter.POST("/login", uc.Login)
	// taskrouter.GET("/get/:id",tc.GetTaskById)
	// taskrouter.PUT("/update/:id",tc.UpdateTask)
	// taskrouter.DELETE("/delete/:id",tc.DeleteTask)
}