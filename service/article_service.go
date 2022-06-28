package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"time"
)

type ArticleService struct {
	BaseService
	ArticleParams
}

func GetArticleService() *ArticleService {
	s := ArticleService{}
	s.model = &s
	return &s
}

func (s *ArticleService) GetTagName() []string {
	return s.TagName
}

func (s *ArticleService) CheckTagName() (int, int) {
	for _, tag_name := range s.TagName {
		if !model.ExistTagByName(tag_name) {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
	}
	return http.StatusOK, e.SUCCESS
}

// func (s *ArticleService) AddArticleTagsByName() (int, int) {
// 	var tags []model.Tag
// 	for _, tag_name := range s.TagName {
// 		tag_id, err := model.GetTagIdByName(tag_name)
// 		if err != nil {
// 			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
// 		}
// 		tag := model.Tag{}
// 		tag.ID = tag_id
// 		tags = append(tags, tag)
// 	}
// 	if err := model.AddArticleTags(s.Id, tags); err != nil {
// 		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
// 	}
// 	return http.StatusOK, e.SUCCESS
// }
func (s *ArticleService) AddArticleTags() (int, int) {
	var tags []model.Tag
	for _, tag_id := range s.TagID {
		if !model.ExistTagByID(tag_id) {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
		tag := model.Tag{}
		tag.ID = uint(tag_id)
		tags = append(tags, tag)
	}
	if err := model.AddArticleTags(s.Id, tags); err != nil {
		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
	}
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) DeleteArticleTags() (int, int) {
	var tags []model.Tag
	for _, tag_id := range s.TagID {
		tag := model.Tag{}
		tag.ID = uint(tag_id)
		tags = append(tags, tag)
	}
	if err := model.DeleteArticleTag(s.Id, tags); err != nil {
		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
	}
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) Add() error {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		if err != nil {
			return err
		}
		tag := model.Tag{}
		tag.ID = tag_id
		tags = append(tags, tag)
	}

	if err := model.AddArticle(
		model.Article{
			OwnerID: s.UserID,
			Content: s.Content,
		}, tags); err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) AddArticleWithImg() error {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		tag := model.Tag{}
		if err != nil {
			tag.Name = tag_name
		} else {
			tag.ID = tag_id
		}
		tags = append(tags, tag)
	}
	var imgs []model.Image
	for _, img_name := range s.ImgName {
		imgs = append(imgs, model.Image{Filename: img_name})
	}
	if err := model.AddArticleWithImg(
		model.Article{
			OwnerID: s.UserID,
			Content: s.Content,
			Title:   s.Title,
		}, tags, imgs); err != nil {
		return err
	}
	return nil
}
func (s *ArticleService) ExistByID() bool {
	return model.ExistArticleByID(s.Id)
}

// func (a *ArticleService) Get() (*model.Article, error) {
// 	var cacheArticle *model.Article

// 	//cache := cache.Article{ID: a.ID}
// 	key := a.GetArticleKey()
// 	if redis.Exists(key) {
// 		data, err := redis.Get(key)
// 		if err != nil {
// 			logging.Info(err)
// 		} else {
// 			json.Unmarshal(data, &cacheArticle)
// 			return cacheArticle, nil
// 		}
// 	}

// 	article, err := model.GetArticle(a.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	redis.Set(key, article, 3600)
// 	return article, nil
// }

func (s *ArticleService) GetArticles(params ArticleGetParams) ([]*model.Article, error) {
	var (
		// articles
		cacheArticles []*model.Article
	)
	key := GetArticleListParamsKey(params.PageNum, params.PageSize)
	if redis.Exists(key) != 0 {
		data := redis.GetBytes(key)
		// if err != nil {
		// 	logging.Info(err)
		// } else {
		json.Unmarshal(data, &cacheArticles)
		// 	return cacheArticles, nil
		// }
		// fmt.Println("ArticleList Cache", data)
		if err := getArticleLikeInfo(cacheArticles, params.Uid); err != nil {
			return nil, err
		}
		return cacheArticles, nil

	}
	// data := make(map[string]interface{})
	// data["uid"] = params.Uid
	articles, err := model.GetArticles(params.PageNum, params.PageSize, nil)
	if err != nil {
		return nil, err
	}
	if err := getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}
	//TODO: 写数据库 缓存一致性问题
	redis.SetBytes(key, articles, time.Minute*3)
	return articles, nil
}

func getArticleLikeInfo(articles []*model.Article, uid int) error {
	for _, article := range articles {
		key := GetModelFieldKey(e.CACHE_ARTICLE, article.ID, e.CACHE_LIKEUSERS)
		if redis.Exists(key) == 0 {
			fmt.Println("SET CACHE_KEY", key)
			likeUsers, err := model.GetArticleLikeUsers(article.ID)
			if err != nil {
				return err
			}
			redis.SetBit(key, 0, 0)
			for _, user := range likeUsers {
				redis.SetBit(key, int64(user.UserId), 1)
			}
		}
		article.LikeCount = redis.BitCount(key)
		// cnt := model.GetArticleLikeCount(article)
		// fmt.Println("LikeCount", cnt)
		if uid != 0 {
			article.IsLike = redis.GetBit(key, int64(uid)) == 1
		}
	}
	return nil
}

func (s *ArticleService) UpdateArticleLike(param ArticleLikeParams) error {
	key := GetModelFieldKey(e.CACHE_ARTICLE, uint(param.Id), e.CACHE_LIKEUSERS)
	messageKey := GetMessageKey(e.CACHE_ARTICLE, uint(param.Id), e.CACHE_LIKEUSERS)
	if redis.Exists(key) != 0 {
		redis.SetBit(key, int64(param.UserID), param.Type)
		redis.SAdd(messageKey, param.UserID)
		return nil
	}
	user := model.User{}
	user.ID = (uint)(param.UserID)
	likeUsers, err := model.GetArticleLikeUsers(uint(param.Id))
	if err != nil {
		return err
	}
	for _, user := range likeUsers {
		redis.SetBit(key, int64(user.UserId), 1)
	}
	redis.SetBit(key, int64(param.UserID), param.Type)
	redis.SAdd(messageKey, param.UserID)
	return nil
	// return model.AddArticleLikeUser(uint(param.Id), user)
}
func (s *ArticleService) Delete() error {
	return model.DeleteArticle(s.Id)
}

func (s *ArticleService) Update() error {
	data := make(map[string]interface{})
	data["state"] = s.State
	data["content"] = s.Content
	return model.UpdateArticle(s.Id, data)
}
func (s *ArticleService) Count(data map[string]interface{}) (int64, error) {
	return model.GetArticleTotal(data)
}

// func (s *ArticleService) Clear() error {
// 	return model.CleanAllArticle()
// }

func (s *ArticleService) Recovery() error {
	return model.RecoverArticle(s.Id)
}

func (s *ArticleService) GetUserID() (uint, error) {
	return model.GetArticleUserID(s.Id)
}

// func (a *ArticleService) GetArticleKey() string {
// 	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
// }
