package model

import (
	"qiu/blog/pkg/e"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	// err := db.Offset(pageNum).Limit(pageSize).Preload(clause.Associations).Find(&articles).Error
	err := db.Offset(pageNum).Limit(pageSize).Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("name", "id")
	}).Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "article_id", "filename")
	}).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	// for _, article := range articles {
	// db.Model(article).Association("LikedUsers").Count()
	// }
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

func GetArticleLikeUsers(id uint) ([]*UserId, error) {
	var like_users []*UserId
	// var like_users []uint
	if err := db.Table(e.TABLE_ARTICLE_LIKE_USERS).Where("`article_id` = ?", id).Select("`user_id`").Find(&like_users).Error; err != nil {
		return nil, err
	}
	// if err := db.Exec("select `user_id` from `blog_article_like_users` where `article_id` = ?", id).Scan(like_users).Error; err != nil {
	// 	return nil, err
	// }
	return like_users, nil
	// if err := tx.Model(&article).Association("LikedUsers").Find(&like_users); err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	// tx.Model(&article).Preload("LikedUsers", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("user_id")
	// }).Find(&like_users)

}
func AddArticleLikeUser(id uint, user User) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	var article Article
	if err := tx.Where("id = ? ", id).First(&article).Error; err != nil {
		return err
	}
	if err := tx.Model(&article).Association("LikedUsers").Append(&user); err != nil {
		tx.Rollback()
		return err
	}
	// cnt := tx.Model(&article).Association("LikedUsers").Count()

	// if err := tx.Model(&article).Update("like_count", cnt).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	return tx.Commit().Error
}

func GetArticleLikeCount(article *Article) int64 {
	return db.Model(article).Association("LikedUsers").Count()
}

func AddArticleLikeUsers(data []ArticleLikeUsers) error {
	// var data []ArticleIdUserId
	// for _, userId := range userIds {
	// 	data = append(data, ArticleIdUserId{ArticleId: articleId, UserId: uint(userId)})
	// }
	return db.Table(e.TABLE_ARTICLE_LIKE_USERS).Clauses(clause.OnConflict{DoNothing: true}).Create(data).Error
}

func DeleteArticleLikeUsers(articleId uint, userIds []uint) error {
	return db.Table(e.TABLE_ARTICLE_LIKE_USERS).Where("article_id = ?", articleId).Where("user_id in ?", userIds).Delete(ArticleIdUserId{}).Error
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
