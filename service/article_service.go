package service

import (
	"encoding/json"
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/logging"
	"qiu/blog/pkg/redis"
	"strconv"
	"strings"
)

type ArticleParams struct {
	Id         uint     `uri:"id"`
	ImgUrl     string   `json:"img_url" form:"img_url"`
	TagName    []string `json:"tag_name" form:"tag_name"`
	Content    string   `json:"content" form:"content"`
	CreatedBy  string   `json:"created_by" form:"created_by"`
	ModifiedBy string   `json:"modified_by" form:"created_by"`
	PageNum    int
	PageSize   int
}
type ArticleService struct {
	BaseService
	ArticleParams
}

func GetArticleService() *ArticleService {
	s := ArticleService{}
	s.model = &s
	return &s
}

// func (s *ArticleService) Bind(c *gin.Context) (int, int) {
// 	var err error

// 	fmt.Println("绑定参数")
// 	if err = c.ShouldBind(s); err != nil {
// 		fmt.Println("绑定错误", err)
// 		return http.StatusBadRequest, e.INVALID_PARAMS
// 	}
// 	// fmt.Println("绑定url")
// 	// if err = c.ShouldBindUri(s); err != nil {
// 	// 	// fmt.Println("绑定数据", s.model)
// 	// 	fmt.Println("绑定错误", err)
// 	// 	return http.StatusBadRequest, e.INVALID_PARAMS
// 	// }

// 	fmt.Println("绑定tag", s.TagName)
// 	fmt.Println("绑定", s.Article)
// 	return http.StatusOK, e.SUCCESS
// }

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

func (s *ArticleService) AddArticleTags() (int, int) {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		if err != nil {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
		tag := model.Tag{}
		tag.ID = tag_id
		tags = append(tags, tag)
	}
	if err := model.AddArticleTags(s.Id, tags); err != nil {
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
			Content:   s.Content,
			CreatedBy: s.CreatedBy,
		}, tags); err != nil {
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

	redis.Set(key, articles, 60)
	return articles, nil
}

func (s *ArticleService) Delete() error {
	return model.DeleteArticle(s.Id)
}

func (s *ArticleService) Count(data map[string]interface{}) (int, error) {
	return model.GetArticleTotal(data)
}

func (s *ArticleService) Clear() error {
	return model.CleanAllArticle()
}

func (s *ArticleService) Recovery() error {
	return model.RecoverArticle(s.Id)
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
