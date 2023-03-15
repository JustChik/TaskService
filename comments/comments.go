package comments

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type (
	CommentService struct {
		db        *sql.DB
		tableName string
	}
	Comment struct {
		Id      uuid.UUID
		Text    string
		Data    time.Time
		Task_id uuid.UUID
	}

	CreateCommentRequest struct {
		Task_id uuid.UUID
		Text    string
	}
	GetCommentRequeset struct {
		Text string
		Data time.Time
	}
)

func NewComment(db *sql.DB, tableName string) *CommentService {
	return &CommentService{
		db:        db,
		tableName: tableName,
	}
}

func (c *CommentService) CreateNewComment(commentCreate CreateCommentRequest) (*Comment, error) {
	id := uuid.New()
	_, err := c.db.Exec(fmt.Sprintf("INSERT INTO %s (id,text,task_id) VALUES ($1,$2,$3);", c.tableName), id.String(), commentCreate.Text, commentCreate.Task_id)
	if err != nil {
		return nil, fmt.Errorf("can't create comment %v", err)
	}
	return nil, nil
}

func (c *CommentService) GetComments() ([]GetCommentRequeset, error) {
	all, err := c.db.Query(fmt.Sprintf("SELECT text, date FROM %s ORDER BY date ASC;", c.tableName))
	if err != nil {
		return nil, fmt.Errorf("can't get commets %v", err)
	}
	com := GetCommentRequeset{}
	res := []GetCommentRequeset{}
	for all.Next() {
		err = all.Scan(&com.Text, &com.Data)
		if err != nil {
			return nil, fmt.Errorf("can't read storage response %v", err)
		}
		res = append(res, com)
	}

	return res, nil
}
