package model

import "gorm.io/gorm"

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
	if err := db.Where("id = ?", replyId).Where("article_id = ?", articleId).First(&comment).Error; err != nil {
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

func GetComments(articleId int, pageNum int, pageSize int) ([]*Comment, error) {
	var comments []*Comment
	err := db.Model(&Comment{}).Where("`article_id` = ?", articleId).Where("`reply_id` IS NULL").Order("created_on desc").Offset(pageNum).Limit(pageSize).
		Select("id", "user_id", "article_id", "created_on", "username", "avator", "content", "like_count").Preload(
		"Replies", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "created_on", "user_id", "article_id", "reply_id", "username", "avator", "content", "like_count")
		}).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
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

func DeleteReply(replyId int) error {
	return db.Delete(&Reply{}, replyId).Error
}

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
