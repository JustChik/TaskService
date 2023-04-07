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
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
		Data time.Time `json:"data"`
	}
	ChangeCommentsRequest struct {
		CommentID uuid.UUID
		NewText   string
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
	now := time.Now()
	_, err := c.db.Exec(fmt.Sprintf("INSERT INTO %s (id,text,task_id, created_at) VALUES ($1,$2,$3, $4);", c.tableName), id.String(), commentCreate.Text, commentCreate.TaskID, now)
	if err != nil {
		return nil, fmt.Errorf("can't create comment %v", err)
	}
	return &Comment{
		Id:     id,
		Text:   commentCreate.Text,
		Data:   now,
		TaskID: commentCreate.TaskID,
	}, nil
}

func (c *Service) GetComments(taskID uuid.UUID) ([]GetCommentRequest, error) {
	all, err := c.db.Query(fmt.Sprintf("SELECT id, text, created_at FROM %s WHERE task_id = '%s' ORDER BY created_at ASC;", c.tableName, taskID))
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

func (c *Service) ChangeComment(commentChange ChangeCommentsRequest) (*Comment, error) {
	now := time.Now()
	_, err := c.db.Exec(fmt.Sprintf("UPDATE %s SET  text=$1, created_at=$2 WHERE id = $3 ;", c.tableName), commentChange.NewText, now, commentChange.CommentID)
	if err != nil {
		return nil, fmt.Errorf("can't change commet %v", err)
	}
	return &Comment{
		Id:   commentChange.CommentID,
		Text: commentChange.NewText,
		Data: now,
	}, nil
}
