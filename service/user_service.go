package service

import (
	"fmt"
	"qiu/blog/model"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	"strconv"
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

func (s *UserService) GetUsernameByID() string {
	return model.GetUsernameByID(s.Id)
}

func (s *UserService) GetUUID(uid uint) string {
	key := strconv.Itoa(int(uid)) + "_" + "uuid"
	uuid := util.GenerateUUID()
	fmt.Println(uuid)
	redis.SetString(key, uuid, 60*60*24)
	return uuid
}

func (s *UserService) CheckUUID(uid uint, uuid string) bool {
	key := strconv.Itoa(int(uid)) + "_" + "uuid"
	if !redis.Exists(key) {
		return false
	}
	v, _ := redis.Get(key)
	fmt.Println("CheckUUID", v, string(v), uuid)
	return uuid == string(v)
}
