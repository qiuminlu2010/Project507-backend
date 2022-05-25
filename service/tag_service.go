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

func (s *TagService) GetCreatedBy() string {
	return s.model.(*model.Tag).CreatedBy
}

func (s *TagService) SetCreatedBy(created_by string) {
	s.model.(*model.Tag).CreatedBy = created_by
}

func (s *TagService) GetModifiedBy() string {
	return s.model.(*model.Tag).ModifiedBy
}

func (s *TagService) SetModifiedBy(modified_by string) {
	s.model.(*model.Tag).ModifiedBy = modified_by
}

func (s *TagService) SetId(id int) {
	s.model.(*model.Tag).ID = id
}
func (s *TagService) Add() error {
	return model.AddTag(s.GetTagModel())
}

func (s *TagService) Delete() bool {
	return model.DeleteTag(s.model.(*model.Tag).ID)
}

func (s *TagService) Update() bool {
	data := make(map[string]interface{})
	data["name"] = s.model.(*model.Tag).Name
	data["modified_by"] = s.model.(*model.Tag).ModifiedBy
	return model.EditTag(s.model.(*model.Tag).ID, data)
}

func (s *TagService) Get(data map[string]int) []model.Tag {
	return model.GetTags(data["pageNum"], data["pageSize"])
}

func (s *TagService) ExistTagByName() bool {
	return model.ExistTagByName(s.model.(*model.Tag).Name)
}
