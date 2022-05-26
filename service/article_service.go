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

type ArticleService struct {
	BaseService
	PageNum  int
	PageSize int
}

func GetArticleService() *ArticleService {
	s := ArticleService{}
	s.model = &model.Article{}
	return &s
}

func (s *ArticleService) AddArticle() error {

	if err := model.AddArticle(*(s.model.(*model.Article))); err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) ExistByID() bool {
	return model.ExistArticleByID(s.model.(*model.Article).ID)
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

func (s *ArticleService) GetArticles(data map[string]interface{}) ([]*model.Article, error) {
	var (
		articles, cacheArticles []*model.Article
	)

	key := s.GetArticlesKey()
	if redis.Exists(key) {
		data, err := redis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := model.GetArticles(s.PageNum, s.PageSize, data)
	if err != nil {
		return nil, err
	}

	redis.Set(key, articles, 3600)
	return articles, nil
}

func (s *ArticleService) Delete() error {
	return model.DeleteArticle(s.model.(*model.Article).ID)
}

func (s *ArticleService) Count(data map[string]interface{}) (int, error) {
	return model.GetArticleTotal(data)
}

// func (a *ArticleService) GetArticleKey() string {
// 	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
// }

func (s *ArticleService) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	// if a.ID > 0 {
	// 	keys = append(keys, strconv.Itoa(a.ID))
	// }
	// if a.TagID > 0 {
	// 	keys = append(keys, strconv.Itoa(a.TagID))
	// }
	// if a.State >= 0 {
	// 	keys = append(keys, strconv.Itoa(a.State))
	// }
	if s.PageNum > 0 {
		keys = append(keys, strconv.Itoa(s.PageNum))
	}
	if s.PageSize > 0 {
		keys = append(keys, strconv.Itoa(s.PageSize))
	}

	return strings.Join(keys, "_")
}
