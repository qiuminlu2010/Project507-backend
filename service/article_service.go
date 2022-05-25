package service

import (
	"encoding/json"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/logging"
	"qiu/blog/pkg/redis"
	"strconv"
	"strings"
)

// type ArticleService struct {
// 	BaseService
// }

// func GetArticleService() *ArticleService {
// 	s := ArticleService{}
// 	s.model = &model.Article{}
// 	return &s
// }
type ArticleService struct {
	model.Article

	PageNum  int
	PageSize int
}

func (a *ArticleService) Add() error {

	if err := model.AddArticle(a.Article); err != nil {
		return err
	}

	return nil
}

func (a *ArticleService) ExistByID() (bool, error) {
	return model.ExistArticleByID(a.ID)
}

func (a *ArticleService) Get() (*model.Article, error) {
	var cacheArticle *model.Article

	//cache := cache.Article{ID: a.ID}
	key := a.GetArticleKey()
	if redis.Exists(key) {
		data, err := redis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := model.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	redis.Set(key, article, 3600)
	return article, nil
}

func (a *ArticleService) GetAll() ([]*model.Article, error) {
	var (
		articles, cacheArticles []*model.Article
	)

	key := a.GetArticlesKey()
	if redis.Exists(key) {
		data, err := redis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := model.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	redis.Set(key, articles, 3600)
	return articles, nil
}

func (a *ArticleService) Delete() error {
	return model.DeleteArticle(a.ID)
}

func (a *ArticleService) Count() (int, error) {
	return model.GetArticleTotal(a.getMaps())
}

func (a *ArticleService) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	// if a.TagID != -1 {
	// 	maps["tag_id"] = a.TagID
	// }

	return maps
}

func (a *ArticleService) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

func (a *ArticleService) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	// if a.TagID > 0 {
	// 	keys = append(keys, strconv.Itoa(a.TagID))
	// }
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}
