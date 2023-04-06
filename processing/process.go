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

	ChangePass struct {
		UserName string
		NewPass  string
		Code     string
	}

	TaskInfo struct {
		UserID      uuid.UUID
		Tittle      string
		Description string
	}

	ChangeTaskRequest struct {
		IDUser      uuid.UUID
		IDTask      uuid.UUID
		Tittle      string
		Description string
	}
	ChangeStatusRequest struct {
		TaskID uuid.UUID
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

	GetTasksRequest struct {
		Tittle      string
		Description string
		Status      tasks.Status
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

func (p *Service) ResetCodeUser(req UserInfo) (string, error) {
	return p.usersSVC.ResetCode(req.UserName)
}

func (p *Service) ChangePasswordUser(req ChangePass) (string, error) {
	return p.usersSVC.ChangePassword(req.UserName, req.Code, req.NewPass)
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
	_, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't change task %v", err)
	}

	return p.tasksSVC.ChangeTask(tasks.ChangeTaskRequest{
		IDUser:      taskChange.IDUser,
		IDTask:      taskChange.IDTask,
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
		IDUser: statusUser.ID,
		IDTask: statusChange.TaskID,
		Status: statusChange.Status,
	})
}
func (p *Service) GetTasks(user UserInfo, req GetTasksRequest) ([]tasks.Task, error) {
	u, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get tasks %v", err)
	}

	filters := map[string]string{
		"user_id": u.ID.String(),
	}
	additionalFilters := req.toMap()
	mergeMaps(filters, additionalFilters)

	return p.tasksSVC.GetTasks(filters)
}

func (p *Service) CreateComment(user UserInfo, createComment CreateCommentRequest) (*comments.Comment, error) {
	commentUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get user %v", err)
	}
	commentTask, err := p.tasksSVC.GetTasks(map[string]string{
		"id":      createComment.TaskID.String(),
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

func (p *Service) GetComments(user UserInfo, getComment GetCommentRequst) ([]comments.GetCommentRequest, error) {
	commentUser, err := p.usersSVC.GetUser(user.UserName, user.UserPass)
	if err != nil {
		return nil, fmt.Errorf("can't get comments %v", err)
	}
	commentTask, err := p.tasksSVC.GetTasks(map[string]string{
		"id":      getComment.TaskID.String(),
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

func (r GetTasksRequest) toMap() map[string]string {
	result := make(map[string]string)
	if r.Tittle != "" {
		result["tittle"] = r.Tittle
	}
	if r.Description != "" {
		result["description"] = r.Description
	}
	if r.Status != "" {
		result["status"] = string(r.Status)
	}

	return result
}

// mergeMaps adds key -> value pairs from new map to old
func mergeMaps(old, new map[string]string) {
	for k, v := range new {
		old[k] = v
	}
}
