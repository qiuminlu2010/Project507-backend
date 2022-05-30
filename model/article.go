package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Article struct {
	Model
	UserID  uint
	Tags    []Tag  `gorm:"many2many:article_tags;"`
	Title   string `json:"title" form:"title"`
	ImgUrl  string `json:"img_url" form:"img_url"`
	Content string `json:"content" form:"content"`
	Like    int    `json:"like" form:"like"`
	Collect int    `json:"collect" form:"collect"`
	Watch   int    `json:"watch" form:"watch"`
	//TODO: Comments   []Comment
	CreatedBy  string `json:"-" form:"created_by" binding:"-"`
	ModifiedBy string `json:"-" form:"created_by" binding:"-"`
	State      int    `json:"-" form:"state" binding:"-"`
}

//通过ID判断文章是否存在
func ExistArticleByID(id uint) bool {
	var art Article
	db.Select("id").Where("id = ?", id).First(&art)
	return art.ID > 0
}

//获取文章数量
func GetArticleTotal(maps interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Article{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

//获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	// err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	err := db.Offset(pageNum).Where(maps).Limit(pageSize).Preload(clause.Associations).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

//通过ID查询文章
func GetArticleById(id uint) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).Preload(clause.Associations).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

//更新文章信息(标签除外)
func UpdateArticle(id uint, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//给文章添加标签
func AddArticleTags(id uint, tags []Tag) error {
	var article Article
	if err := db.Where("id = ? ", id).First(&article).Error; err != nil {
		return err
	}
	if err := db.Model(&article).Association("Tags").Append(tags); err != nil {
		return err
	}
	return nil
}

//给文章删除标签
func DeleteArticleTag(id uint, tags []Tag) error {
	var article Article
	if err := db.Where("id = ? ", id).First(&article).Error; err != nil {
		return err
	}
	if err := db.Model(&article).Association("Tags").Delete(tags); err != nil {
		return err
	}
	return nil
}

// AddArticle add a single article
func AddArticle(article Article, tags []Tag) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&article).Association("Tags").Append(tags); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeleteArticle delete a single article
func DeleteArticle(id uint) error {
	var article Article
	article.ID = id

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// if err := tx.Model(&article).Association("Tags").Clear().Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	if err := tx.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func RecoverArticle(id uint) error {
	var article Article
	article.ID = id
	if err := db.Unscoped().Model(&article).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

func CleanAllArticle() error {
	if err := db.Unscoped().Select(clause.Associations).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}

func GetArticleUserID(id uint) (uint, error) {
	var article Article
	if err := db.Select("user_id").Where("id = ?", id).First(&article).Error; err != nil {
		return 0, err
	}
	return article.UserID, nil
}
