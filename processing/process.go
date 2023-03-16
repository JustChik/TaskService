package processing

import (
	"TaskService/comments"
	"TaskService/tasks"
	"TaskService/users"
)

type (
	Service struct {
		usersSVC    *users.Service
		tasksSVC    *tasks.Service
		commentsSVC *comments.Service
	}
)

func NewService(usersSVC *users.Service, tasksSVC *tasks.Service, commentsSVC *comments.Service) *Service {
	return &Service{
		usersSVC:    usersSVC,
		tasksSVC:    tasksSVC,
		commentsSVC: commentsSVC,
	}
}

//func (p *Service) CreateNewUser(username, pass string) (*users.User, error) {
//	res, err := p.usersSVC.Create(username, pass)
//	return
//}
