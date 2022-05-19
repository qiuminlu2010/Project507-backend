package service

import (
	"fmt"
	"qiu/blog/model"
)

type UserService struct {
	BaseService
}

func GetUserService() *UserService {
	s := UserService{}
	s.model = &model.User{}
	return &s
}
func (s *UserService) Register() error {
	return model.AddUser(s.GetUserModel())
}

func (s *UserService) Delete() bool {
	return model.DeleteUser(s.model.(*model.User).ID)
}

func (s *UserService) Login() (bool, error) {
	return model.ValidLogin(s.GetUsername(), s.GetPassword())
}

func (s *UserService) IfExisted() bool {
	return model.ExistUsername(s.GetUsername())
}

func (s *UserService) UpdatePassword() bool {
	data := make(map[string]interface{})
	data["password"] = s.GetPassword()
	return model.UpdatePassword(s.model.(*model.User).ID, data)
}

func (s *UserService) GetPassword() string {
	fmt.Println("密码", s.model.(*model.User).Password)
	return s.model.(*model.User).Password
}

func (s *UserService) GetUsername() string {
	return s.model.(*model.User).Username
}

func (s *UserService) GetUsernameByID() string {
	return model.GetUsernameByID(s.model.(*model.User).ID)
}

func (s *UserService) GetUserModel() model.User {
	return *(s.model.(*model.User))
}
