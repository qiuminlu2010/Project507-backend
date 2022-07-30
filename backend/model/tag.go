package model

import (
	"qiu/backend/pkg/e"

	"gorm.io/gorm/clause"
)

//https://blog.csdn.net/weixin_45604257/article/details/105139862

// func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
// 	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

// 	return
// }
func GetTags() ([]*TagInfo, error) {
	var tags []*TagInfo
	if err := db.Model(&Tag{}).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagsByPrefix(prefix string, count int) ([]*TagInfo, error) {
	var tags []*TagInfo
	if err := db.Model(&Tag{}).Where("name like ?", prefix+"%").Limit(count).Order("name").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// func GetTagArticles(pageNum int, pageSize int, tagName string) ([]int, error) {
// 	var articles []*Article
// 	tagIdsql := db.Model(&Tag{}).Where("name = ?", tagName).Select("id")
// 	articleIdSql := db.Table(e.TABLE_ARTICLE_TAGS).Where("tag_id = ?", tagIdsql).Select("article_id")
// 	err:=db.Model(&Article{}).Where("id in (?)",articleIdSql).Offset(pageNum).Limit(pageSize).Order("created_on desc").Find(&articles).Error
// 	if err := db.Table(e.TABLE_ARTICLE_TAGS).Where(
// 		"tag_id = ?", db.Model(&Tag{}).Where("name = ?", tagName).
// 			Select("id")).Select("article_id").Find(&articleIds).Error; err != nil {
// 		return nil, err
// 	}
// 	return articleIds, nil
// }
func GetTagArticleIds(tagName string) ([]int, error) {
	// db.Model(&Tag{}).Where("name = ?",tagName).Select("id")
	var articleIds []int
	if err := db.Table(e.TABLE_ARTICLE_TAGS).Where(
		"tag_id = (?)", db.Model(&Tag{}).Where("name = ?", tagName).
			Select("id")).Select("article_id").Find(&articleIds).Error; err != nil {
		return nil, err
	}
	// err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	// var tag Tag
	// tag.ID = tag_id
	// err = db.Model(&tag).Association("Articles").Find(&articles)
	return articleIds, nil
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
func DeleteTag(id int) error {
	var tag Tag
	tag.ID = uint(id)
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
