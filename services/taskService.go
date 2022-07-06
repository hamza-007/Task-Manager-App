package services

import (
	"database/sql"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hamza-007/Task-Manager-App/models"
)


type TaskService interface {
	Add(*models.Task,string) error
	GetTasks(string)([]models.Task,error)
	GetTask(string,string)(models.Task, error)
	UpdateTask(*models.Task,string,string) (models.Task,error)
	DeleteTask(string,string) (models.Task,error)
}

type TaskSvc struct{
	BD 	*sql.DB
}

func NewTaskService(bd *sql.DB) TaskService {
	return &TaskSvc{
		BD: bd,
	}
}

func (ts *TaskSvc)Add(t *models.Task,id string) error {
	stmt, err := ts.BD.Prepare("INSERT INTO tasks(id,description,completed,created_at,completed_at,userid) VALUES (?,?,?,?,?,?) ")
	if err != nil {
		return err 
	}
	token, err := jwt.ParseWithClaims(id, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		
		return err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	_, err = stmt.Exec(&t.Guid, &t.Description, &t.Completed, &t.CreatedAt,&t.CompletedAt,&claims.Issuer)
	return err
}

func (ts *TaskSvc)GetTasks(id string) ([]models.Task, error) {
	var tasks []models.Task
	var task models.Task
	token, err := jwt.ParseWithClaims(id, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		
		return nil,err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	rows, e := ts.BD.Query("SELECT * FROM tasks WHERE userid = ?",&claims.Issuer)
	
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&task.Guid, &task.Description, &task.Completed, &task.CreatedAt,&task.CompletedAt,&task.Userid)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (ts *TaskSvc)GetTask(id string,userid string) (models.Task, error) {
	var task models.Task
	token, err := jwt.ParseWithClaims(userid, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		
		return task,err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	row, err := ts.BD.Query("SELECT * FROM tasks WHERE id = ? AND userid = ?", id,claims.Issuer)
	if err != nil {
		return task, err
	}
	for row.Next() {
		err := row.Scan(&task.Guid, &task.Description, &task.Completed, &task.CreatedAt,&task.CompletedAt,&task.Userid)
		if err != nil {
			return task, err
		}
	}
	return task, nil
}

func (ts *TaskSvc)UpdateTask(t *models.Task,id string,userid string) (models.Task,error) {
	var task models.Task
	token, err := jwt.ParseWithClaims(userid, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		
		return task,err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	task , err = ts.GetTask(id,userid)
	if err != nil { return task,err }
	stmt, err := ts.BD.Prepare("UPDATE tasks SET description = ? , completed = ? , completed_at = ? WHERE id = ? AND userid = ? ")
	if err != nil {
		return task,err
	}
	_, err = stmt.Exec( &t.Description, &t.Completed, time.Now().Format(time.ANSIC) ,&id,&claims.Issuer)
	return task,err 
}


func (ts *TaskSvc)DeleteTask(id string,userid string) (models.Task,error) {
	task , err := ts.GetTask(id,userid)
	if err != nil {
		return task,err
	}
	token, err := jwt.ParseWithClaims(userid, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		
		return task,err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	stmt, err := ts.BD.Prepare("DELETE FROM tasks WHERE id = ? AND userid = ?")

	if err != nil { return task,err }

	_, err = stmt.Exec(id,claims.Issuer)
	return task,err 
}