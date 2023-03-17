package processing

import (
	"TaskService/comments"
	"TaskService/tasks"
	"TaskService/users"
	"fmt"
	"github.com/google/uuid"
)

type (
	Service struct {
		usersSVC    *users.Service
		tasksSVC    *tasks.Service
		commentsSVC *comments.Service
	}
	UserInfo struct {
		UserName string
		UserPass string
	}

	TaskInfo struct {
		UserID      uuid.UUID
		Tittle      string
		Description string
	}

	ChangeTaskRequest struct {
		Id          uuid.UUID
		Tittle      string
		Description string
	}
	ChangeStatusRequest struct {
		Status tasks.Status
	}
	ChekTask struct {
		TaskID uuid.UUID
	}
	CreateCommentRequest struct {
		TaskID uuid.UUID
		Text   string
	}
	GetCommentRequst struct {
		TaskID uuid.UUID
	}
)

func NewService(usersSVC *users.Service, tasksSVC *tasks.Service, commentsSVC *comments.Service) *Service {
	return &Service{
		usersSVC:    usersSVC,
		tasksSVC:    tasksSVC,
		commentsSVC: commentsSVC,
	}
}

func (p *Service) CreateUser(req UserInfo) (*users.User, error) {
	return p.usersSVC.Create(req.UserName, req.UserPass)
}

func (p *Service) GetUser(req UserInfo) (*users.User, error) {
	return p.usersSVC.GetUser(req.UserName, req.UserPass)
}

func (p *Service) GetAllUsers() ([]users.User, error) {
	return p.usersSVC.GetAllUsers()
}

func (p *Service) CreateTask(user UserInfo, taskInfo TaskInfo) (*tasks.Task, error) {
	userInfo, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't create task %v", err)
	}

	return p.tasksSVC.CreateNewTask(tasks.CreateTaskRequest{
		UserID:      userInfo.ID,
		Tittle:      taskInfo.Tittle,
		Description: taskInfo.Description,
	})

}

func (p *Service) ChangeTask(user UserInfo, taskChange ChangeTaskRequest) (*tasks.Task, error) {
	taskUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't change task %v", err)
	}
	return p.tasksSVC.ChangeTask(tasks.ChangeTaskRequest{
		Id:          taskUser.ID,
		Tittle:      taskChange.Tittle,
		Description: taskChange.Description,
	})
}

func (p *Service) ChangeTaskStatus(user UserInfo, statusChange ChangeStatusRequest) (tasks.Status, error) {
	statusUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return "", fmt.Errorf("can't change task %v", err)
	}
	return p.tasksSVC.ChangeStatus(tasks.ChangeStatusRequest{
		Id:     statusUser.ID,
		Status: statusChange.Status,
	})
}
func (p *Service) GetTasks(user UserInfo, filters map[string]string) ([]tasks.Task, error) {
	_, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get tasks %v", err)
	}
	return p.tasksSVC.GetTasks(filters)
}

func (p *Service) CreateComment(user UserInfo, task ChekTask, createComment CreateCommentRequest) (*comments.Comment, error) {
	commentUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get user %v", err)
	}
	commentTask, err := p.tasksSVC.GetTasks(map[string]string{
		"id":      task.TaskID.String(),
		"user_id": commentUser.ID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("can't create comment %v", err)
	}

	if len(commentTask) == 0 {
		return nil, fmt.Errorf("task does not exist")
	}

	createComment.TaskID = commentTask[0].Id
	return p.commentsSVC.CreateNewComment(comments.CreateCommentRequest{
		TaskID: createComment.TaskID,
		Text:   createComment.Text,
	})
}

func (p *Service) GetComments(user UserInfo, task ChekTask, getComment GetCommentRequst) ([]comments.GetCommentRequest, error) {
	commentUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get comments %v", err)
	}
	commentTask, err := p.tasksSVC.GetTasks(map[string]string{
		"id":      task.TaskID.String(),
		"user_id": commentUser.ID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("can't get comment %v", err)
	}

	if len(commentTask) == 0 {
		return nil, fmt.Errorf("task does not exist")
	}

	getComment.TaskID = commentTask[0].Id
	return p.commentsSVC.GetComments(getComment.TaskID)
}
