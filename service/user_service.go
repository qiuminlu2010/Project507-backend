package service

import (
	"qiu/blog/model"
)

type UserParams struct {
	Id       uint   `uri:"id"`
	Username string `json:"username" form:"username" binding:"omitempty,printascii,gte=3,lte=20"`
	Password string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=20"`
}

type UserService struct {
	BaseService
	UserParams
}

func GetUserService() *UserService {
	s := UserService{}
	s.model = &s
	return &s
}

func (s *UserService) Register() error {
	return model.AddUser(model.User{
		Username: s.Username,
		Password: s.Password,
	})
}

func (s *UserService) Delete() bool {
	return model.DeleteUser(s.Id)
}

func (s *UserService) Login() (bool, error) {
	return model.ValidLogin(s.Username, s.Password)
}

func (s *UserService) IfExisted() bool {
	return model.ExistUsername(s.Username)
}

func (s *UserService) UpdatePassword() bool {
	data := make(map[string]interface{})
	data["password"] = s.Password
	return model.UpdatePassword(s.Id, data)
}

func (s *UserService) GetUsernameByID() string {
	return model.GetUsernameByID(s.Id)
}
