package service

import (
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
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

func (s *TagService) Get() []*model.TagInfo {
	tags, err := model.GetTags()
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		redis.ZAdd(e.CACHE_TAGS, 0, tag.Name)
	}
	return tags
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

func setTagCache() {
	tags, err := model.GetTags()
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		redis.ZAdd(e.CACHE_TAGS, 0, tag.Name)
	}
}

func (s *TagService) Hint(params *TagsGetParams) ([]*model.TagInfo, error) {
	// if redis.Exists(e.CACHE_TAGS) == 0 {
	// 	setTagCache()
	// }
	// start := params.TagName
	// end := start + "龟"
	// for i:=1;i<=int(redis.ZCard(e.CACHE_TAGS));i++ {
	// 	if(redis.ZRandMember())
	// }
	var tags []*model.TagInfo
	var err error
	tags, err = model.GetTagsByPrefix(params.TagName, params.PageSize)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

//TODO: 基于redis的自动补全 1.插入和删除有序集合在并发环境下不适用，需要上锁 2.顺序查找 3.二分查找有序集合
func (s *TagService) HintByCache(params *TagsGetParams) []string {
	if redis.Exists(e.CACHE_TAGS) == 0 {
		setTagCache()
	}
	var tags []string
	start := params.TagName
	end := start + "龟"
	var i int64
	cnt := 0
	for i = 1; i <= redis.ZCard(e.CACHE_TAGS); i++ {
		tag := redis.ZRange(e.CACHE_TAGS, i-1, i-1)[0]
		if tag >= start && tag < end {
			tags = append(tags, tag)
			cnt++
		}
		if tag == end || cnt == params.PageSize {
			return tags
		}
	}
	return tags
}
