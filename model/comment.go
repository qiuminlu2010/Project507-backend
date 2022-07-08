package model

import (
	"fmt"
	"qiu/blog/pkg/e"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddComment(userId int, articleId int, content string) error {
	comment := Comment{UserID: uint(userId), ArticleID: uint(articleId), Content: content}
	if err := db.Model(&User{}).Where("id = ?", userId).Select("username", "avator").First(&comment).Error; err != nil {
		return err
	}
	if err := db.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func AddReply(userId int, articleId int, replyId int, content string) error {
	var comment Comment
	// comment.ID = uint(replyId)
	// comment.ArticleID = uint(articleId)
	if err := db.Where("id = ?", replyId).Where("article_id = ?", articleId).Where("`reply_id` IS NULL").First(&comment).Error; err != nil {
		return err
	}
	reply := Comment{UserID: uint(userId), ArticleID: uint(articleId), Content: content}
	if err := db.Model(&User{}).Where("id = ?", userId).Select("username", "avator").First(&reply).Error; err != nil {
		return err
	}
	if err := db.Model(&comment).Association("Replies").Append(&reply); err != nil {
		return err
	}
	return nil
}

func GetComments(articleId int, userId int, pageNum int, pageSize int) ([]*Comment, error) {
	var comments []*Comment
	likeCountSql := ",(select count(*) from `blog_user_like_comments` where `blog_comment`.`id` = comment_id) as like_count"
	selectSql := "`id`,`user_id`,`article_id`,`created_on`,`username`,`avator`,`content`" + likeCountSql
	selectReplySql := "`id`,`user_id`,`article_id`,`reply_id`,`created_on`,`username`,`avator`,`content`" + likeCountSql
	isLikeSql := ""
	if userId > 0 {
		isLikeSql = fmt.Sprintf(",(select count(*) from `blog_user_like_comments` where `blog_comment`.`id` = comment_id and user_id = %d) as is_like", userId)
		selectSql += isLikeSql
		selectReplySql += isLikeSql
	}
	err := db.Table("blog_comment").
		Where("`article_id` = ?", articleId).
		Where("`reply_id` IS NULL").
		Order("created_on desc").
		Offset(pageNum).
		Limit(pageSize).
		Select(selectSql).
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Select(selectReplySql)
		}).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// isLikeSql := ""
	// if userId > 0 {
	// 	isLikeSql = fmt.Sprintf(",(select count(*) from `blog_user_like_comments` where `blog_comment`.`id` = comment_id and user_id = %d) as is_like", userId)
	// 	selectSql += isLikeSql
	// }
	return comments, nil
}

func AddCommentLike(userId int, CommentId int) error {
	// data := make(map[string]interface{})
	// data["user_id"] = userId
	// data["comment_id"] = CommentId
	data := CommentIdUserId{UserId: uint(userId), CommentID: uint(CommentId)}
	return db.Table(e.TABLE_USER_LIKE_COMMENTS).Clauses(clause.OnConflict{DoNothing: true}).Create(&data).Error
}

func DeleteCommentLike(userId int, CommentId int) error {
	return db.Table(e.TABLE_USER_LIKE_COMMENTS).Where("comment_id = ?", CommentId).Where("user_id = ?", userId).Delete(CommentIdUserId{}).Error
}

func DeleteComment(commentId int) error {
	return db.Delete(&Comment{}, commentId).Error
}

// func AddReply(userId int, commentId int, content string) error {
// 	reply := Reply{UserID: uint(userId), CommentID: uint(commentId), Content: content}
// 	if err := db.Model(&User{}).Where("id = ?", userId).Select("username", "avator").First(&reply).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Create(&reply).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DeleteReply(replyId int) error {
// 	return db.Delete(&Reply{}, replyId).Error
// }

func GetArticleOwnerIdByCommentId(commentId int) (uint, error) {
	var article Article
	if err := db.Select("owner_id").Where("id = (?)", db.Model(&Comment{}).Select("article_id").Where("id = ?", commentId)).First(&article).Error; err != nil {
		return 0, err
	}
	return article.OwnerID, nil
}

func GetArticleOwnerIdByReplyId(replyId int) (uint, error) {
	var article Article
	q1 := db.Model(&Reply{}).Select("comment_id").Where("id = ?", replyId)
	q2 := db.Model(&Comment{}).Select("article_id").Where("id = (?)", q1)
	if err := db.Select("owner_id").Where("id = (?)", q2).First(&article).Error; err != nil {
		return 0, err
	}
	return article.OwnerID, nil
}
