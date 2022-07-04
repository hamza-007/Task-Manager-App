package services

import (
	"github.com/hamza-007/Task-Manager-App/models"
	"database/sql"
	"time"
)


type TaskService interface {
	Add(*models.Task) error
	GetTasks()([]models.Task,error)
	GetTask(string)(models.Task, error)
	UpdateTask(*models.Task,string) (models.Task,error)
	DeleteTask(string) (models.Task,error)
}

type TaskSvc struct{
	BD 	*sql.DB
}

func NewTaskService(bd *sql.DB) TaskService {
	return &TaskSvc{
		BD: bd,
	}
}

func (ts *TaskSvc)Add(t *models.Task) error {
	stmt, err := ts.BD.Prepare("INSERT INTO tasks(id,description,completed,created_at,updated_at) VALUES (?,?,?,?,?) ")
	if err != nil {
		return err 
	}
	_, err = stmt.Exec(&t.Guid, &t.Description, &t.Completed, &t.CreatedAt,&t.UpdatedAt)
	return err
}

func (ts *TaskSvc)GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	var task models.Task
	rows, e := ts.BD.Query("SELECT * FROM tasks")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&task.Guid, &task.Description, &task.Completed, &task.CreatedAt,&task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (ts *TaskSvc)GetTask(id string) (models.Task, error) {
	var task models.Task
	row, err := ts.BD.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		return task, err
	}
	for row.Next() {
		err := row.Scan(&task.Guid, &task.Description, &task.Completed, &task.CreatedAt,&task.UpdatedAt)
		if err != nil {
			return task, err
		}
	}
	return task, nil
}

func (ts *TaskSvc)UpdateTask(t *models.Task,id string) (models.Task,error) {
	task , err := ts.GetTask(id)
	if err != nil { return task,err }
	stmt, err := ts.BD.Prepare("UPDATE tasks SET description = ? , completed = ? , updated_at = ? WHERE id = ? ")
	if err != nil {
		return task,err
	}
	_, err = stmt.Exec( &t.Description, &t.Completed, time.Now().Format(time.ANSIC) ,&id)
	return task,err 
}


func (ts *TaskSvc)DeleteTask(id string) (models.Task,error) {
	task , err := ts.GetTask(id)
	if err != nil {
		return task,err
	}

	stmt, err := ts.BD.Prepare("DELETE FROM tasks WHERE id = ? ")

	if err != nil { return task,err }

	_, err = stmt.Exec(id)

	return task,err 
}