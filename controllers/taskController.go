package controllers

import (
	"database/sql"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hamza-007/Task-Manager-App/handlers"
	"github.com/hamza-007/Task-Manager-App/models"
	"github.com/hamza-007/Task-Manager-App/services"
)

type TaskController struct {
	TaskService services.TaskService
}

func NewTaskController(taskservice services.TaskService) TaskController {
	return TaskController{
		TaskService: taskservice,
	}
}

func (tc *TaskController) AddTask(c *gin.Context){
	var task models.Task
	if err := c.ShouldBindJSON(&task) ; err != nil{
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	task.Guid = uuid.New().String()
	task.CreatedAt = time.Now().Format(time.ANSIC)

	cookie ,err:= c.Cookie("jwt")
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	
	if err := tc.TaskService.Add(&task,cookie) ; err != nil{
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	var res = handlers.NewHTTPResponse(http.StatusCreated,task)
	c.JSON(http.StatusCreated,res)
}


func (tc *TaskController) GetAllTasks(c *gin.Context){
	cookie ,err:= c.Cookie("jwt")
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	tasks,err := tc.TaskService.GetTasks(cookie)
	if err != nil{
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	} 
	if len(tasks) == 0 {
		var res = handlers.NewHTTPResponse(http.StatusNotFound, sql.ErrNoRows)
		c.JSON(http.StatusNotFound, res)
		return
	}
	var res = handlers.NewHTTPResponse(http.StatusOK, tasks)
	c.JSON(http.StatusOK, res)
}


func (tc *TaskController) GetTaskById(c *gin.Context){
	cookie ,err:= c.Cookie("jwt")
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	id := c.Param("id")
	task ,err := tc.TaskService.GetTask(id,cookie)
	if err != nil  {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if task.Guid == "" {
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, "task doesn't exist !! ")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	var res = handlers.NewHTTPResponse(http.StatusOK, task)
	c.JSON(http.StatusOK, res)
}



func (tc *TaskController) UpdateTask(c *gin.Context){
	cookie ,err:= c.Cookie("jwt")
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	id := c.Param("id")
	var t models.Task
	if err := c.ShouldBindJSON(&t) ; err != nil{
		var res = handlers.NewHTTPResponse(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	task , err := tc.TaskService.UpdateTask(&t,id,cookie)
	if err != nil  {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	} 
	if task.Guid == "" {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, "task doesn't exist !! ")
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	var res = handlers.NewHTTPResponse(http.StatusOK, "task updated succesfully !!")
	c.JSON(http.StatusOK, res)
}


func (tc *TaskController) DeleteTask(c *gin.Context){
	cookie ,err:= c.Cookie("jwt")
	if err!=nil {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	id := c.Param("id")
	task ,err := tc.TaskService.DeleteTask(id,cookie) 
	if err != nil{
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, err)
		c.JSON(http.StatusInternalServerError, res)
		return
	} 
	if task.Guid == "" {
		var res = handlers.NewHTTPResponse(http.StatusInternalServerError, "task doesn't exist")
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	var res = handlers.NewHTTPResponse(http.StatusOK, "task deleted succesfully !!")
	c.JSON(http.StatusOK, res)
}


func (tc *TaskController) RegisterTaskRoutes(rg *gin.RouterGroup) {
	taskrouter := rg.Group("/task")
	taskrouter.POST("/create", tc.AddTask)
	taskrouter.GET("/get", tc.GetAllTasks)
	taskrouter.GET("/get/:id",tc.GetTaskById)
	taskrouter.PUT("/update/:id",tc.UpdateTask)
	taskrouter.DELETE("/delete/:id",tc.DeleteTask)
}