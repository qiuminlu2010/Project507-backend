package tag_service

import (
	"qiu/backend/model"
	"qiu/backend/pkg/e"
	"qiu/backend/pkg/redis"
	base "qiu/backend/service/base"
	param "qiu/backend/service/param"
)

type TagService struct {
	base.BaseService
}

var tagService TagService

func GetTagService() *TagService {
	return &tagService
}

func (s *TagService) Add(name string) error {
	return model.AddTag(model.Tag{
		Name: name,
	})
}

func (s *TagService) Delete(id int) error {
	return model.DeleteTag(id)
}

func (s *TagService) Get(params *param.PageGetParams) []*model.TagInfo {
	tags, err := model.GetTags()
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		redis.ZAdd(e.CACHE_TAGS, 0, tag.Name)
	}
	return tags
}

func (s *TagService) Recovery(tagId int) error {
	return model.RecoverTag(uint(tagId))
}

func (s *TagService) Clear(tagId int) error {
	return model.CleanTag(uint(tagId))
}

func (s *TagService) ExistTagByName(name string) bool {
	return model.ExistTagByName(name)
}

func (s *TagService) Hint(params *param.TagsGetParams) ([]*model.TagInfo, error) {
	var tags []*model.TagInfo
	var err error
	tags, err = model.GetTagsByPrefix(params.TagName, params.PageSize)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

//TODO: 基于redis的自动补全 1.插入和删除有序集合在并发环境下不适用，需要上锁 2.顺序查找边界 3.二分查找有序集合的关键词边界
func (s *TagService) HintByCache(params *param.TagsGetParams) []string {
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

func setTagCache() {
	tags, err := model.GetTags()
	if err != nil {
		panic(err)
	}
	for _, tag := range tags {
		redis.ZAdd(e.CACHE_TAGS, 0, tag.Name)
	}
}
