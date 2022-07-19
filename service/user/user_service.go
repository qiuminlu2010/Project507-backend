package user_service

import (
	"encoding/json"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	base "qiu/blog/service/base"
	cache "qiu/blog/service/cache"
	param "qiu/blog/service/param"
)

type UserService struct {
	base.BaseService
}

var userService UserService

func GetUserService() *UserService {
	return &userService
}

func (s *UserService) CountUser(data map[string]interface{}) (int64, error) {
	return model.GetUserTotal(data)
}

func (s *UserService) GetUserList(param *param.UserListGetParams) ([]*model.User, error) {
	return model.GetUserList(param.PageNum, param.PageSize)
}

func (s *UserService) GetUserInfo(userId int) (*model.UserInfo, error) {
	userInfo, err := model.GetUserInfo(uint(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (s *UserService) Add(params *param.UserAddParams) error {
	var user model.User
	user.Username = params.Username
	user.Password = params.Password
	if params.Name == "" {
		user.Name = params.Username
	} else {
		user.Name = params.Name
	}
	return model.AddUser(&user)
}

func (s *UserService) Delete(id int) error {
	return model.DeleteUser(uint(id))
}

func (s *UserService) Update(params *param.UserUpdateParams) error {
	data := make(map[string]interface{})
	if params.Name != "" {
		data["name"] = params.Name
	}
	if params.Password != "" {
		data["password"] = params.Password
	}
	return model.UpdateUser(uint(params.UserId), data)
}
func (s *UserService) Login(params *param.UserLoginParams) (*model.UserInfo, error) {
	return model.ValidLogin(params.Username, params.Password)
}

func (s *UserService) ExistUsername(name string) bool {
	return model.ExistUsername(name) == 1
}

// func (s *UserService) UpdatePassword() error {
// 	data := make(map[string]interface{})
// 	data["password"] = s.Password
// 	return model.UpdateUser(s.Id, data)
// }

func (s *UserService) UpdateState(params *param.UserUpdateParams) error {
	data := make(map[string]interface{})
	data["state"] = params.State
	return model.UpdateUser(uint(params.UserId), data)
}
func (s *UserService) GetUsernameByID(id int) string {
	return model.GetUsernameByID(uint(id))
}

func (s *UserService) GetUUID(uid uint) string {
	key := cache.GetModelFieldKey("user", uid, "uuid")
	uuid := util.GenerateUUID()
	redis.Set(key, uuid, 60*60*24)
	return uuid
}

func (s *UserService) CheckUUID(uid uint, uuid string) bool {
	key := cache.GetModelFieldKey("user", uid, "uuid")
	if redis.Exists(key) == 0 {
		return false
	}
	v := redis.Get(key)
	return uuid == v
}

func (s *UserService) GetUsersByName(params *param.UsersGetParams) ([]*model.UserBase, error) {
	return model.GetUsersByName(params.Name+"%", params.PageNum, params.PageSize)
}

//关注列表： 1.设置缓存 user:id:follows 2.设置缓存 user:id
func (s *UserService) GetFollows(params *param.FollowsGetParams) ([]*model.UserBase, error) {
	//TODO: 分页
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	var followIds []int
	var err error

	if redis.Exists(key) == 0 {
		if err := setUserFollowCache(params.UserId); err != nil {
			return nil, err
		}
	}

	value := redis.SGET(key)
	followIds, err = util.StringsToInts(value)

	if err != nil {
		return nil, err
	}

	followUsers := GetUsersCache(followIds)
	return followUsers, nil
}

func GetUsersCache(userIds []int) []*model.UserBase {
	var userInfos []*model.UserBase
	for _, userId := range userIds {
		userKey := cache.GetModelIdKey(e.CACHE_USER, userId)
		if redis.Exists(userKey) == 0 {
			err := setUserCache(userId)
			if err != nil {
				continue
			}
		}
		var userInfo model.UserBase
		json.Unmarshal(redis.GetBytes(userKey), &userInfo)
		redis.Expire(userKey, e.DURATION_USER_INFO)
		userInfos = append(userInfos, &userInfo)
	}
	return userInfos

}

func setUserCache(userId int) error {
	userInfo, err := model.GetUser(uint(userId))
	key := cache.GetModelIdKey(e.CACHE_USER, userId)
	if err != nil {
		return err
	}
	redis.SetBytes(key, userInfo, e.DURATION_FOLLOWS)
	// redis.Expire(key,e.DURATION_FOLLOWS)
	return nil
}

func setUserFollowCache(userId int) error {
	//user:id:follows
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_FOLLOWS)
	followIds, err := model.GetFollowIds(uint(userId))
	if err != nil {
		return err
	}
	for _, followId := range followIds {
		redis.SAdd(key, followId)
	}
	redis.Expire(key, e.DURATION_FOLLOWS)
	return nil
}

//关注操作： 1.设置缓存 user:id:follows 2.修改缓存 3 添加消息缓存 message:user:id:follows
func (s *UserService) UpsertFollowUser(params *param.UpsertUserFollowParams) error {

	key := cache.GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	// messageKey := GetMessageKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)

	// if redis.Exists(key) == 0 {
	// 	if err := setUserFollowCache(params.UserId); err != nil {
	// 		return err
	// 	}
	// }

	// m := make(map[string]interface{})
	// m[strconv.Itoa(params.FollowId)] = params.Type

	// redis.HashSet(messageKey, m)

	if params.Type == 1 {
		model.FollowUser(params.UserId, params.FollowId)
		if redis.Exists(key) != 0 {
			redis.SAdd(key, params.FollowId)
		}

	} else {
		model.UnFollowUser(params.UserId, params.FollowId)
		if redis.Exists(key) != 0 {
			redis.SDEL(key, params.FollowId)
		}
	}

	return nil
}

func (s *UserService) GetFans(userId int) ([]*model.UserBase, error) {
	fanIds, err := model.GetFanIds(uint(userId))
	if err != nil {
		return nil, err
	}
	return GetUsersCache(fanIds), nil
}

func GetUserCache(userId int) *model.UserBase {
	var userInfo model.UserBase

	userKey := cache.GetModelIdKey(e.CACHE_USER, userId)
	if redis.Exists(userKey) == 0 {
		err := setUserCache(userId)
		if err != nil {
			return nil
		}
	}
	json.Unmarshal(redis.GetBytes(userKey), &userInfo)
	redis.Expire(userKey, e.DURATION_USER_INFO)

	return &userInfo

}
