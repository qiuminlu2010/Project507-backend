package service

import (
	"qiu/blog/model"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
)

type UserParams struct {
	Id       uint   `uri:"id"`
	Username string `json:"username" form:"username" binding:"omitempty,printascii,gte=3,lte=20"`
	Password string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
	State    int    `json:"state" form:"state"`
}

type UserService struct {
	BaseService
	UserParams
	PageNum  int
	PageSize int
}

func GetUserService() *UserService {
	s := UserService{}
	s.model = &s
	return &s
}

func (s *UserService) CountUser(data map[string]interface{}) (int64, error) {
	return model.GetUserTotal(data)
}

func (s *UserService) GetUserList(data map[string]interface{}) ([]*model.User, error) {
	return model.GetUserList(s.PageNum, s.PageSize, data)
}

func (s *UserService) Add() error {
	return model.AddUser(model.User{
		Username: s.Username,
		Password: s.Password,
	})
}

func (s *UserService) Delete() error {
	return model.DeleteUser(s.Id)
}

func (s *UserService) Login() (model.User, error) {
	return model.ValidLogin(s.Username, s.Password)
}

func (s *UserService) ExistUsername() error {
	return model.ExistUsername(s.Username)
}

func (s *UserService) UpdatePassword() error {
	data := make(map[string]interface{})
	data["password"] = s.Password
	return model.UpdateUser(s.Id, data)
}

func (s *UserService) UpdateState() error {
	data := make(map[string]interface{})
	data["state"] = s.State
	return model.UpdateUser(s.Id, data)
}
func (s *UserService) GetUsernameByID() string {
	return model.GetUsernameByID(s.Id)
}

func (s *UserService) GetUUID(uid uint) string {
	key := GetKeyName("user", uid, "uuid")
	uuid := util.GenerateUUID()
	redis.Set(key, uuid, 60*60*24)
	return uuid
}

func (s *UserService) CheckUUID(uid uint, uuid string) bool {
	key := GetKeyName("user", uid, "uuid")
	if redis.Exists(key) == 0 {
		return false
	}
	v := redis.Get(key)
	return uuid == v
}
