package service

import (
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	"strconv"
)

type UserService struct {
	BaseService
	UserLoginParams
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
	key := GetModelKey("user", uid, "uuid")
	uuid := util.GenerateUUID()
	redis.Set(key, uuid, 60*60*24)
	return uuid
}

func (s *UserService) CheckUUID(uid uint, uuid string) bool {
	key := GetModelKey("user", uid, "uuid")
	if redis.Exists(key) == 0 {
		return false
	}
	v := redis.Get(key)
	return uuid == v
}

func (s *UserService) UpsertFollowUser(params UpsertUserFollowParams) error {

	key := GetModelKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	messageKey := GetMessageKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)

	if redis.Exists(key) != 0 {
		// redis.SetBit(key, int64(params.UserId), params.Type)
		m := make(map[string]interface{})
		m[strconv.Itoa(params.FollowId)] = params.Type

		redis.HashSet(messageKey, m)

		if params.Type == 1 {
			redis.SAdd(key, params.FollowId)
		} else {
			redis.SDEL(key, params.FollowId)
		}

		return nil
	}

	Follows, err := model.GetFollows(uint(params.UserId))

	if err != nil {
		return err
	}

	for _, follow := range Follows {
		m := make(map[string]interface{})
		m[strconv.Itoa(follow.FollowId)] = 1
		redis.SAdd(key, m)
	}

	m := make(map[string]interface{})
	m[strconv.Itoa(params.FollowId)] = params.Type

	redis.HashSet(messageKey, m)

	if params.Type == 1 {
		redis.SAdd(key, params.FollowId)
	} else {
		redis.SDEL(key, params.FollowId)
	}

	return nil
}
