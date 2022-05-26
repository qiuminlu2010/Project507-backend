package model

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	Model
	Tags    []Tag  `gorm:"many2many:article_tags;"`
	Title   string `json:"title" form:"title"`
	ImgUrl  string `json:"img_url" form:"img_url"`
	Content string `json:"content" form:"content"`
	Like    int    `json:"like" form:"like"`
	Collect int    `json:"collect" form:"collect"`
	//TODO: Comments   []Comment
	CreatedBy  string `json:"created_by" form:"created_by"`
	ModifiedBy string `json:"modified_by" form:"created_by"`
	State      int    `json:"state"`
}

//通过ID判断文章是否存在
func ExistArticleByID(id int) bool {
	var art Article
	db.Select("id").Where("id = ?", id).First(&art)
	return art.ID > 0
}

//获取文章数量
func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

//获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	// err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	err := db.Offset(pageNum).Limit(pageSize).Association("Tag").Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

//通过ID查询文章
func GetArticleById(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&article).Association("Tag").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

//更新文章信息(标签除外)
func UpdateArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//给文章添加标签
func AddArticleTags(id int, tags []Tag) error {
	err := db.Model(&Article{}).Association("Tag").Append(tags).Error
	if err != nil {
		return err
	}
	return nil
}

// AddArticle add a single article
func AddArticle(article Article) error {
	if err := db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

// DeleteArticle delete a single article
func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}
	return nil
}

// CleanAllArticle clear all article
func CleanAllArticle() error {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}
