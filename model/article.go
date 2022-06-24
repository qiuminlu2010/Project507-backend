package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Image struct {
	Model
	ArticleID uint   `json:"-" form:"-" binding:"-"`
	Filename  string `json:"filename" form:"filename" binding:"-"`
}
type Article struct {
	Model
	OwnerID   uint    `json:"owner_id"`
	User      User    `gorm:"foreignkey:OwnerID" binding:"-" json:"-"` // 使用 UserRefer 作为外键
	Tags      []Tag   `gorm:"many2many:article_tags;" json:"tags"`
	Images    []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"images"`
	Title     string  `json:"title" form:"title"`
	Content   string  `json:"content" form:"content"`
	LikeCount int     `json:"like_count" form:"like_count" binding:"-"`
	// Collect int     `json:"collect" form:"collect" binding:"-"`
	// Watch   int     `json:"watch" form:"watch" binding:"-"`
	LikedUsers []User `gorm:"many2many:article_like_users;" json:"-"`
	IsLike     bool   `json:"is_like" form:"is_like" binding:"-"`
	//TODO: Comments   []Comment
	CreatedBy  string `json:"-" form:"created_by" binding:"-"`
	ModifiedBy string `json:"-" form:"created_by" binding:"-"`
	State      int    `json:"state" form:"state" binding:"-"`
}

type Comment struct {
	Model
	UserID  uint   `json:"user_id"`
	Content string `json:"content" form:"content"`
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
	err := db.Where("id = ?", id).Preload(clause.Associations).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

//更新文章信息(标签除外)
func UpdateArticle(id uint, data interface{}) error {
	return db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
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

// AddArticle add a single article
func AddArticleWithImg(article Article, tags []Tag, imgs []Image) error {
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
	if err := tx.Model(&article).Association("Images").Append(imgs); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func AddArticleLikeUser(id uint, user User) error {
	var article Article
	if err := db.Where("id = ? ", id).First(&article).Error; err != nil {
		return err
	}
	if err := db.Model(&article).Association("LikedUsers").Append(&user); err != nil {
		return err
	}
	return nil
}

// DeleteArticle delete a single article
func DeleteArticle(id uint) error {
	var article Article
	article.ID = id
	return db.Where("id = ?", id).Delete(Article{}).Error
}

func RecoverArticle(id uint) error {
	var article Article
	article.ID = id
	return db.Unscoped().Model(&article).Update("deleted_at", nil).Error
}

// func CleanAllArticle() error {
// 	if err := db.Unscoped().Select(clause.Associations).Delete(&Article{}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetArticleUserID(id uint) (uint, error) {
	var article Article
	if err := db.Select("user_id").Where("id = ?", id).First(&article).Error; err != nil {
		return 0, err
	}
	return article.OwnerID, nil
}
