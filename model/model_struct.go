package model

import (
	"gorm.io/gorm"
)

type Model struct {
	ID         uint           `gorm:"primary_key" uri:"id" `
	CreatedOn  int            `gorm:"index" binding:"-" json:"created_on,omitempty"`
	ModifiedOn int            `binding:"-" json:"modified_on,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index"  binding:"-" json:"-"`
}

type Image struct {
	ID         uint   `gorm:"primary_key" uri:"id" `
	CreatedOn  int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
	ModifiedOn int    `binding:"-" json:"modified_on,omitempty"`
	ArticleID  uint   `json:"-" form:"-" binding:"-"`
	Filename   string `json:"filename" form:"filename" binding:"-"`
	// Thumbnail int    `json:"-" form:"-"`
}
type Article struct {
	Model
	OwnerID   uint   `gorm:"index" json:"owner_id"`
	User      User   `gorm:"foreignkey:OwnerID" binding:"-" json:"-"` // 使用 OwnerID  作为外键
	Title     string `json:"title" form:"title"`
	Content   string `json:"content" form:"content"`
	LikeCount int64  `json:"like_count" form:"like_count" binding:"-"`
	// Collect int     `json:"collect" form:"collect" binding:"-"`
	// Watch   int     `json:"watch" form:"watch" binding:"-"`
	State      int     `json:"state" form:"state" binding:"-"`
	IsLike     bool    `json:"is_like" form:"is_like" binding:"-"`
	Tags       []Tag   `gorm:"many2many:article_tags;" json:"tags"`
	Images     []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignkey:ArticleID" json:"images"`
	LikedUsers []User  `gorm:"many2many:article_like_users;" json:"-"`
	//TODO: Comments   []Comment
	// CreatedBy  string `json:"-" form:"created_by" binding:"-"`
	// ModifiedBy string `json:"-" form:"created_by" binding:"-"`

}

type ArticleInfo struct {
	ID        uint    `json:"id" form:"id"`
	CreatedOn int     `binding:"-" json:"created_on,omitempty"`
	OwnerID   uint    `json:"owner_id"`
	Title     string  `json:"title" form:"title"`
	Content   string  `json:"content" form:"content"`
	LikeCount int64   `json:"like_count" form:"like_count" binding:"-"`
	IsLike    bool    `json:"is_like" form:"is_like" binding:"-"`
	Images    []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignkey:ArticleID" json:"images"`
	Tags      []Tag   `gorm:"many2many:article_tags;" json:"tags"`
}

type ArticleCache struct {
	ID        uint `json:"id" form:"id"`
	CreatedOn int  `binding:"-" json:"created_on,omitempty"`
}

type Comment struct {
	Model
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	ReplyID   *uint     `json:"reply_id"`
	Username  string    `json:"username"`
	Avator    string    `json:"avator"`
	Content   string    `json:"content" form:"content" binding:"gte=1,lte=200"`
	LikeCount int       `json:"like_count" `
	IsLike    int       `json:"is_like"`
	Replies   []Comment `gorm:"foreignkey:ReplyID" json:"replies"`
}

type Reply struct {
	Model
	UserID    uint   `json:"user_id"`
	CommentID uint   `json:"comment_id"`
	Username  string `json:"username"`
	Avator    string `json:"avator"`
	Content   string `json:"content" form:"content" binding:"gte=1,lte=200"`
	LikeCount int    `json:"like_count" `
}

type CommentIdUserId struct {
	CommentID uint `json:"comment_id"`
	UserId    uint `json:"user_id"`
}

// type CommentInfo struct {
// 	ID        uint    `json:"id" form:"id"`
// 	CreatedOn int     `gorm:"index" binding:"-" json:"created_on,omitempty"`
// 	UserID    uint    `json:"user_id"`
// 	ArticleID uint    `json:"article_id"`
// 	Username  string  `json:"username"`
// 	Avator    string  `json:"avator"`
// 	Content   string  `json:"content" form:"content" binding:"gte=1,lte=200"`
// 	Replies   []Reply `json:"replies"`
// 	LikeCount int     `json:"like_count" `
// }
// type ReplyInfo struct {
// 	ID        uint   `json:"id" form:"id"`
// 	CreatedOn int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
// 	UserID    uint   `json:"user_id"`
// 	CommentID uint   `json:"comment_id"`
// 	Username  string `json:"username"`
// 	Avator    string `json:"avator"`
// 	Content   string `json:"content" form:"content" binding:"gte=1,lte=200"`
// 	LikeCount int    `json:"like_count" `
// }
type User struct {
	Model
	Username  string `json:"username" form:"username" binding:"omitempty,printascii,gte=6,lte=20" gorm:"index;unique"`
	Password  string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
	StudentId string `json:"student_id" form:"student_id" binding:"omitempty,numeric"`
	// Articles     []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"articles,omitempty" binding:"-"`
	// FollowNum    int        `json:"follow_num" form:"follow_num" binding:"-"`
	// FanNum       int        `json:"fan_num" form:"fan_num" binding:"-"`
	Avator       string     `json:"avator" form:"avator"`
	State        int        `json:"state" form:"state" binding:"gte=0,lte=1"`
	LikeArticles []*Article `gorm:"many2many:article_like_users" binding:"-" json:"like_articles"`
	Follows      []*User    `gorm:"many2many:user_follows"`
	LikeComments []*Comment `gorm:"many2many:user_like_comments" binding:"-" json:"-"`
}

type ArticleLikeUsers struct {
	ArticleID int
	UserID    int
	CreatedAt int `gorm:"index"  binding:"-" json:"created_at,omitempty"`
}

type UserArticlesCache struct {
	ID        uint
	OwnerID   int
	CreatedOn int
}

type UserBase struct {
	ID       uint   `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Avator   string `json:"avator" form:"avator"`
	// FollowNum int    `json:"follow_num" form:"follow_num" binding:"-"`
	// FanNum    int    `json:"fan_num" form:"fan_num" binding:"-"`
	// TODO: LikeNum int
}

type UserInfo struct {
	UserBase
	FollowNum int64 `json:"follow_num" form:"follow_num" binding:"-"`
	FanNum    int64 `json:"fan_num" form:"fan_num" binding:"-"`
	// TODO: LikeNum int
}
type UserId struct {
	UserId uint `json:"user_id"`
}

type FollowId struct {
	Id int `json:"follow_id"`
}

type ArticleIdTagId struct {
	ArticleId uint `json:"article_id"`
	TagId     uint `json:"tag_id"`
}
type ArticleIdUserId struct {
	ArticleId uint `json:"article_id"`
	UserId    uint `json:"user_id"`
}

type UserIdFollowId struct {
	UserId   uint `json:"user_id"`
	FollowId uint `json:"follow_id"`
}
type Tag struct {
	ID        uint   `gorm:"primary_key" uri:"id" `
	CreatedOn int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
	Name      string `json:"name" form:"name" binding:"required,lte=20" gorm:"index;unique;not null"`
	//Type       string `json:"type" form:"type" `
	// CreatedBy  string    `json:"-" form:"created_by" binding:"-" `
	// ModifiedBy string    `json:"-" form:"modified_by" binding:"-" `
	State    int       `json:"-" form:"state" binding:"-" `
	Articles []Article `gorm:"many2many:article_tags;" binding:"-" json:"-"`
}

type TagInfo struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
