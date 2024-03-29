package tasks

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

const (
	StatusOpen       Status = "Open"
	StatusInProgress Status = "In Progress"
	StatusOnHold     Status = "On Hold"
	StatusDone       Status = "Done"
	StatusCanceled   Status = "Canceled"
)

type (
	Service struct {
		db        *sql.DB
		tableName string
	}
	Task struct {
		Id          uuid.UUID `json:"id"`
		Tittle      string    `json:"tittle"`
		Description string    `json:"description"`
		Status      Status    `json:"status"`
	}
	CreateTaskRequest struct {
		UserID      uuid.UUID
		Tittle      string
		Description string
	}

	Status string

	ChangeTaskRequest struct {
		IDUser      uuid.UUID
		IDTask      uuid.UUID
		Tittle      string
		Description string
	}
	ChangeStatusRequest struct {
		IDUser uuid.UUID
		IDTask uuid.UUID
		Status Status
	}
)

func NewService(db *sql.DB, tableName string) *Service {
	return &Service{
		db:        db,
		tableName: tableName,
	}
}

func (t *Service) CreateNewTask(taskCreate CreateTaskRequest) (*Task, error) {
	id := uuid.New()
	status := StatusOpen
	_, err := t.db.Exec(fmt.Sprintf("INSERT INTO %s (id,tittle,description,status,user_id) VALUES ($1,$2,$3,$4,$5);", t.tableName), id.String(), taskCreate.Tittle, taskCreate.Description, status, taskCreate.UserID)
	if err != nil {
		return nil, fmt.Errorf("can't create new taskCreate %v", err)
	}
	return &Task{
		Id:          id,
		Tittle:      taskCreate.Tittle,
		Description: taskCreate.Description,
		Status:      StatusOpen,
	}, nil
}

func (t *Service) ChangeTask(taskChange ChangeTaskRequest) (*Task, error) {
	_, err := t.db.Exec(fmt.Sprintf("UPDATE %s SET tittle = $1, description = $2 WHERE id = $3 AND user_id = $4;", t.tableName), taskChange.Tittle, taskChange.Description, taskChange.IDTask, taskChange.IDUser)
	if err != nil {
		return nil, fmt.Errorf("can't change task %v", err)
	}
	return &Task{
		Id:          taskChange.IDTask,
		Tittle:      taskChange.Tittle,
		Description: taskChange.Description,
	}, nil
}

func (t *Service) ChangeStatus(changeStatus ChangeStatusRequest) (Status, error) {
	_, err := t.db.Exec(fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 AND user_id = $3;", t.tableName), changeStatus.Status, changeStatus.IDTask, changeStatus.IDUser)
	if err != nil {
		return "", fmt.Errorf("can't change Status %v", err)
	}
	return changeStatus.Status, nil
}

func (t *Service) GetTasks(filters map[string]string) ([]Task, error) {
	query := fmt.Sprintf("SELECT id,tittle,description,status FROM %s WHERE 1=1", t.tableName)
	if len(filters) > 0 {
		for k, val := range filters {
			query += " AND "
			if k == "user_id" || k == "id" {
				query += fmt.Sprintf("%s = '%s' ", k, val)
			} else {
				query += fmt.Sprintf("%s LIKE '%%%s%%'", k, val)
			}

		}
	}

	query += ";"
	res, err := t.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("can't find tasks %v", err)
	}

	u := Task{}
	all := []Task{}
	for res.Next() {
		err := res.Scan(&u.Id, &u.Tittle, &u.Description, &u.Status)
		if err != nil {
			return nil, fmt.Errorf("can't read storage response %v", err)
		}
		all = append(all, u)
	}
	return all, nil

}
