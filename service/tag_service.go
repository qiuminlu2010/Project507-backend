package service

import (
	// "fmt"
	"qiu/blog/model"
)

type TagService struct {
	BaseService
}

func GetTagService() *TagService {
	s := TagService{}
	s.model = &model.Tag{}
	return &s
}

func (s *TagService) GetTagModel() model.Tag {
	return *(s.model.(*model.Tag))
}

func (s *TagService) GetTagCreatedBy() string {
	return s.model.(*model.Tag).CreatedBy
}

func (s *TagService) Add() error {
	return model.AddTag(s.GetTagModel())
}

func (s *TagService) Delete() bool {
	return model.DeleteTag(s.model.(*model.Tag).ID)
}

func (s *TagService) Update() bool {
	data := make(map[string]interface{})
	data["password"] = s.model.(*model.Tag).Name
	return model.EditTag(s.model.(*model.User).ID, data)
}

func (s *TagService) Get(data map[string]int) []model.Tag {
	return model.GetTags(data["pageNum"], data["pageSize"])
}

func (s *TagService) ExistTagByName() bool {
	return model.ExistTagByName(s.model.(*model.Tag).Name)
}
