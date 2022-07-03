package service

import (
	"qiu/blog/model"
	"qiu/blog/pkg/e"
)

type TagParams struct {
	Id         uint   `uri:"id"`
	Name       string `json:"name" form:"name" binding:"omitempty,lte=20,gte=2" `
	CreatedBy  string `json:"created_by" form:"created_by" `
	ModifiedBy string `json:"modified_by" form:"modified_by"`
	PageNum    int
	PageSize   int
}

type TagService struct {
	BaseService
	TagParams
}

func GetTagService() *TagService {
	s := TagService{}
	s.model = &s
	return &s
}

func (s *TagService) GetCreatedBy() string {
	return s.CreatedBy
}

func (s *TagService) SetCreatedBy(created_by string) {
	s.CreatedBy = created_by
}

func (s *TagService) GetModifiedBy() string {
	return s.ModifiedBy
}

func (s *TagService) SetModifiedBy(modified_by string) {
	s.ModifiedBy = modified_by
}

func (s *TagService) Add() error {
	return model.AddTag(model.Tag{
		Name: s.Name,
		// CreatedBy: s.CreatedBy,
	})
}

func (s *TagService) Delete() error {
	return model.DeleteTag(s.Id)
}

func (s *TagService) Update() error {
	data := make(map[string]interface{})
	data["name"] = s.Name
	data["modified_by"] = s.ModifiedBy
	return model.EditTag(s.Id, data)
}

func (s *TagService) Get() []model.Tag {
	return model.GetTags(s.PageNum, s.PageSize)
}
func (s *TagService) GetArticles(params *TagArticleGetParams) ([]*model.ArticleInfo, error) {
	articleIds, err := model.GetTagArticleIds(params.TagName)
	if err != nil {
		return nil, err
	}
	articles, err := getArticlesCache(articleIds, e.CACHE_ARTICLES)
	if err != nil {
		return nil, err
	}
	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}
	return articles, nil
	// articles, err := model.GetArticlesById(params.PageNum, params.PageSize, articleIds)
	// if err != nil {
	// 	return nil, err
	// }
	// return articles, nil
}
func (s *TagService) Recovery() error {
	return model.RecoverTag(s.Id)
}
func (s *TagService) Clear() error {
	return model.CleanTag(s.Id)
}
func (s *TagService) ExistTag() bool {
	return model.ExistTagByName(s.Name)
}

func (s *TagService) ExistTagByName(name string) bool {
	return model.ExistTagByName(name)
}
