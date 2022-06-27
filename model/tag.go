package model

import "gorm.io/gorm/clause"

//https://blog.csdn.net/weixin_45604257/article/details/105139862

// func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
// 	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

// 	return
// }
func GetTags(pageNum int, pageSize int) (tags []Tag) {
	db.Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagArticles(tag_id uint) (articles []Article, err error) {
	// err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	var tag Tag
	tag.ID = tag_id
	err = db.Model(&tag).Association("Articles").Find(&articles)
	return
}

func GetTagTotal(maps interface{}) (count int64) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	return tag.ID > 0
}

func AddTag(tag Tag) error {
	return db.Create(&tag).Error
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	return tag.ID > 0
}

func GetTagIdByName(name string) (uint, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil {
		return 0, err
	}
	return tag.ID, nil
}
func DeleteTag(id uint) error {
	var tag Tag
	tag.ID = id
	return db.Where("id = ?", id).Delete(&Tag{}).Error
}

func RecoverTag(id uint) error {
	var tag Tag
	tag.ID = id
	if err := db.Unscoped().Model(&tag).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

func EditTag(id uint, data interface{}) error {
	return db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
}

func CleanTag(id uint) error {
	var tag Tag
	tag.ID = id
	return db.Unscoped().Select(clause.Associations).Delete(&tag).Error
}
