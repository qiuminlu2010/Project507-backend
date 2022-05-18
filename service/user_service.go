package service

import (
	"qiu/blog/model"
)

type UserService struct {
	model.User
}

func (user *UserService) Register() error {
	return model.AddUser(user.User)
}

func (user *UserService) Login() (bool, error) {
	return model.ValidLogin(user.Username, user.Password)
}

func (user *UserService) IfExisted() bool {
	return model.ExistUsername(user.Username)
}
