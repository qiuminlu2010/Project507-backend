package service

import (
	"encoding/json"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	"strconv"
	"time"
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
	key := GetModelFieldKey("user", uid, "uuid")
	uuid := util.GenerateUUID()
	redis.Set(key, uuid, 60*60*24)
	return uuid
}

func (s *UserService) CheckUUID(uid uint, uuid string) bool {
	key := GetModelFieldKey("user", uid, "uuid")
	if redis.Exists(key) == 0 {
		return false
	}
	v := redis.Get(key)
	return uuid == v
}

func (s *UserService) GetFollows(params UserFollowsParams) ([]model.UserInfo, error) {

	key := GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	var followIds []int
	var err error

	if redis.Exists(key) == 0 {
		if err := setUserFollowCache(params.UserId); err != nil {
			return nil, err
		}
		// followIds, err = model.GetFollowIds(uint(params.UserId))
		// if err != nil {
		// 	return nil, err
		// }
		// for _, followId := range followIds {
		// 	redis.SAdd(key, followId)
		// }
	}

	value := redis.SGET(key)
	followIds, err = util.StringsToInts(value)

	if err != nil {
		return nil, err
	}

	followUsers := getUserCache(followIds)
	return followUsers, nil
}

func getUserCache(userIds []int) []model.UserInfo {
	var userInfos []model.UserInfo
	for _, userId := range userIds {
		userKey := GetModelIdKey(e.CACHE_USER, userId)
		if redis.Exists(userKey) == 0 {
			err := setUserCache(userId)
			if err != nil {
				continue
			}
		}
		var userInfo model.UserInfo
		json.Unmarshal(redis.GetBytes(userKey), &userInfo)
		userInfos = append(userInfos, userInfo)
	}
	return userInfos

}

func setUserCache(userId int) error {
	userInfo, err := model.GetUser(uint(userId))
	key := GetModelIdKey(e.CACHE_USER, userId)
	if err != nil {
		return err
	}
	redis.SetBytes(key, userInfo, 6*time.Hour)
	return nil
}

func setUserFollowCache(userId int) error {
	key := GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_FOLLOWS)
	followIds, err := model.GetFollowIds(uint(userId))
	if err != nil {
		return err
	}
	for _, followId := range followIds {
		redis.SAdd(key, followId)
	}
	return nil
}

func (s *UserService) UpsertFollowUser(params UpsertUserFollowParams) error {

	key := GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	messageKey := GetMessageKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)

	if redis.Exists(key) == 0 {
		if err := setUserFollowCache(params.UserId); err != nil {
			return err
		}
	}
	// if redis.Exists(key) != 0 {

	// 	m := make(map[string]interface{})
	// 	m[strconv.Itoa(params.FollowId)] = params.Type

	// 	redis.HashSet(messageKey, m)

	// 	if params.Type == 1 {
	// 		redis.SAdd(key, params.FollowId)
	// 	} else {
	// 		redis.SDEL(key, params.FollowId)
	// 	}

	// 	return nil
	// }

	// Follows, err := model.GetFollows(uint(params.UserId))

	// if err != nil {
	// 	return err
	// }

	// for _, follow := range Follows {
	// 	// m := make(map[string]interface{})
	// 	// m[strconv.Itoa(follow.FollowId)] = 1
	// 	redis.SAdd(key, follow.ID)
	// }

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
