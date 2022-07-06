package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/hamza-007/Task-Manager-App/controllers"
	"github.com/hamza-007/Task-Manager-App/db"
	"github.com/hamza-007/Task-Manager-App/services"
	"github.com/joho/godotenv"
)

var (
	tc controllers.TaskController
	ts services.TaskService
	uc controllers.UserController
	us services.UserService
	bd *sql.DB
)


func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalln("error loading env file")
	}
	router := gin.Default()
	bd = db.Connect()
	ts = services.NewTaskService(bd)
	tc = controllers.NewTaskController(ts)
	us = services.NewUserService(bd)
	uc = controllers.NewUserController(us)


	mainpath := os.Getenv("MAIN__PATH")
	tc.RegisterTaskRoutes(router.Group(mainpath))
	uc.RegisterUserRoutes(router.Group(mainpath))
	
	router.Use(func(c *gin.Context){
		c.String(http.StatusBadRequest," Invalid Request !! !")
	})

	defer bd.Close()
	defer log.Fatal(http.StatusAccepted)
	log.Fatalln(router.Run(os.Getenv("PORT")))
	
}
