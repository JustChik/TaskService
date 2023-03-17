package comments

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type (
	Service struct {
		db        *sql.DB
		tableName string
	}
	Comment struct {
		Id     uuid.UUID
		Text   string
		Data   time.Time
		TaskID uuid.UUID
	}

	CreateCommentRequest struct {
		TaskID uuid.UUID
		Text   string
	}

	GetCommentRequest struct {
		ID   uuid.UUID
		Text string
		Data time.Time
	}
)

func NewService(db *sql.DB, tableName string) *Service {
	return &Service{
		db:        db,
		tableName: tableName,
	}
}

func (c *Service) CreateNewComment(commentCreate CreateCommentRequest) (*Comment, error) {
	id := uuid.New()
	_, err := c.db.Exec(fmt.Sprintf("INSERT INTO %s (id,text,task_id) VALUES ($1,$2,$3);", c.tableName), id.String(), commentCreate.Text, commentCreate.TaskID)
	if err != nil {
		return nil, fmt.Errorf("can't create comment %v", err)
	}
	return nil, nil
}

func (c *Service) GetComments(taskID uuid.UUID) ([]GetCommentRequest, error) {
	all, err := c.db.Query(fmt.Sprintf("SELECT id, text, date FROM %s WHERE task_id = '%s' ORDER BY date ASC;", c.tableName, taskID))
	if err != nil {
		return nil, fmt.Errorf("can't get commets %v", err)
	}
	com := GetCommentRequest{}
	res := []GetCommentRequest{}
	for all.Next() {
		err = all.Scan(&com.ID, &com.Text, &com.Data)
		if err != nil {
			return nil, fmt.Errorf("can't read storage response %v", err)
		}
		res = append(res, com)
	}

	return res, nil
}
