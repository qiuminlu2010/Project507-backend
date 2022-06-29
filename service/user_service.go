package service

import (
	"encoding/json"
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

func (s *UserService) GetUserInfo(userId int) (*model.UserInfo, error) {
	userInfo, err := model.GetUserInfo(uint(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
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

//关注列表： 1.设置缓存 user:id:follows 2.设置缓存 user:id
func (s *UserService) GetFollows(params *UserGetParams) ([]*model.UserBase, error) {
	//TODO: 分页
	key := GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
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

	followUsers := getUserCache(followIds)
	return followUsers, nil
}

func getUserCache(userIds []int) []*model.UserBase {
	var userInfos []*model.UserBase
	for _, userId := range userIds {
		userKey := GetModelIdKey(e.CACHE_USER, userId)
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
	key := GetModelIdKey(e.CACHE_USER, userId)
	if err != nil {
		return err
	}
	redis.SetBytes(key, userInfo, e.DURATION_FOLLOWS)
	// redis.Expire(key,e.DURATION_FOLLOWS)
	return nil
}

func setUserFollowCache(userId int) error {
	//user:id:follows
	key := GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_FOLLOWS)
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
func (s *UserService) UpsertFollowUser(params UpsertUserFollowParams) error {

	key := GetModelFieldKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)
	messageKey := GetMessageKey(e.CACHE_USER, uint(params.UserId), e.CACHE_FOLLOWS)

	if redis.Exists(key) == 0 {
		if err := setUserFollowCache(params.UserId); err != nil {
			return err
		}
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

func (s *UserService) GetFans(userId int) ([]*model.UserBase, error) {
	fanIds, err := model.GetFanIds(uint(userId))
	if err != nil {
		return nil, err
	}
	return getUserCache(fanIds), nil
}

//(userId, pageNum, pageSize) => []ArticleInfo
func (s *UserService) GetUserArticles(params *ArticleGetParams) ([]*model.ArticleInfo, error) {
	//user:id:articles
	key := GetModelFieldKey(e.CACHE_USER, uint(params.Uid), e.CACHE_ARTICLES)
	var articles []*model.ArticleInfo
	var err error
	if redis.Exists(key) == 0 {
		if err = setUserArticleCache(params.Uid); err != nil {
			return nil, err
		}
	}
	//有序集合，根据时间戳降排
	articleIds := redis.ZRevRange(key, int64(params.PageNum), int64(params.PageSize-1))
	articles, err = getArticlesCache(articleIds, key)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

//设置缓存： user:id:articles (score:time member:articleId)
func setUserArticleCache(userId int) error {
	key := GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_ARTICLES)
	articles, err := model.GetUserArticles(uint(userId))
	if err != nil {
		return err
	}
	for _, article := range articles {
		redis.ZAdd(key, float64(article.CreatedOn), article.ID)
	}
	redis.Expire(key, e.DURATION_USERARTICLES)
	return nil
}
