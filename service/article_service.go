package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/logging"
	"qiu/blog/pkg/redis"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArticleService struct {
	BaseService
	model.Article
	PageNum  int
	PageSize int
	TagName  []string `json:"tag_name" form:"tag_name"`
}

func GetArticleService() *ArticleService {
	s := ArticleService{}
	return &s
}
func (s *ArticleService) Bind(c *gin.Context) (int, int) {
	var err error

	fmt.Println("绑定参数")
	if err = c.ShouldBind(s); err != nil {
		fmt.Println("绑定错误", err)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	// fmt.Println("绑定url")
	// if err = c.ShouldBindUri(s); err != nil {
	// 	// fmt.Println("绑定数据", s.model)
	// 	fmt.Println("绑定错误", err)
	// 	return http.StatusBadRequest, e.INVALID_PARAMS
	// }

	fmt.Println("绑定tag", s.TagName)
	fmt.Println("绑定", s.Article)
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) GetTagName() []string {
	return s.TagName
}

func (s *ArticleService) AddTags() (int, int) {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		_, err := model.GetTagIdByName(tag_name)
		if err != nil {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
		// tag := model.Tag{}
		// tag.ID = tag_id
		tags = append(tags, model.Tag{Name: tag_name})
	}

	if err := model.AddArticleTags(s.Article, tags); err != nil {
		fmt.Println("添加文章标签失败", err)
		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
	}

	return http.StatusOK, e.SUCCESS
}
func (s *ArticleService) Add() error {

	if err := model.AddArticle(s.Article); err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) ExistByID() bool {
	return model.ExistArticleByID(s.ID)
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
	return model.DeleteArticle(s.ID)
}

func (s *ArticleService) Count(data map[string]interface{}) (int, error) {
	return model.GetArticleTotal(data)
}

func (s *ArticleService) GetCreatedBy() string {
	return s.CreatedBy
}

func (s *ArticleService) SetCreatedBy(created_by string) {
	s.CreatedBy = created_by
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
