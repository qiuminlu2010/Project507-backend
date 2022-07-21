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

type User struct {
	Model
	Username     string     `json:"username" form:"username" binding:"printascii,gte=6,lte=20" gorm:"size:255;index;unique"`
	Password     string     `json:"password" form:"password" binding:"printascii,gte=6,lte=100" gorm:"size:255"`
	StudentId    string     `json:"student_id" form:"student_id" binding:"omitempty,numeric"`
	Name         string     `json:"name" form:"name" binding:"printascii,gte=6,lte=20" gorm:"size:255;index;unique"`
	Avatar       string     `json:"avatar" form:"avatar"`
	State        int        `json:"state" form:"state" binding:"gte=0,lte=1"`
	LikeArticles []*Article `gorm:"many2many:article_like_users" binding:"-" json:"like_articles"`
	Follows      []*User    `gorm:"many2many:user_follows"`
	LikeComments []*Comment `gorm:"many2many:user_like_comments" binding:"-" json:"-"`
}

type Article struct {
	Model
	OwnerID   uint   `gorm:"index" json:"owner_id"`
	User      User   `gorm:"foreignkey:OwnerID" binding:"-" json:"-"` // 使用 OwnerID  作为外键
	Title     string `json:"title" form:"title" gorm:"size:255;"`
	Content   string `gorm:"collate:utf8" json:"content" form:"content" binding:"gte=1,lte=500"`
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
	ID            uint    `json:"id" form:"id"`
	CreatedOn     int     `binding:"-" json:"created_on,omitempty"`
	OwnerID       uint    `json:"owner_id"`
	OwnerName     string  `json:"owner_name"`
	OwnerUsername string  `json:"owner_username"`
	OwnerAvatar   string  `json:"owner_avatar"`
	Title         string  `json:"title" form:"title"`
	Content       string  `gorm:"collate:utf8" json:"content" form:"content"`
	LikeCount     int64   `json:"like_count" form:"like_count" binding:"-"`
	IsLike        bool    `json:"is_like" form:"is_like" binding:"-"`
	Images        []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignkey:ArticleID" json:"images"`
	Tags          []Tag   `gorm:"many2many:article_tags;" json:"tags"`
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
	Avatar    string    `json:"avatar"`
	Content   string    `gorm:"collate:utf8" json:"content" form:"content" binding:"gte=1,lte=500"`
	LikeCount int       `json:"like_count" `
	IsLike    int       `json:"is_like"`
	Replies   []Comment `gorm:"foreignkey:ReplyID" json:"replies"`
}

type Reply struct {
	Model
	UserID    uint   `json:"user_id"`
	CommentID uint   `json:"comment_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content" form:"content" binding:"gte=1,lte=500"`
	LikeCount int    `json:"like_count" `
}

type CommentIdUserId struct {
	CommentID uint `json:"comment_id"`
	UserId    uint `json:"user_id"`
}

type ArticleLikeUsers struct {
	ArticleID int `gorm:"primaryKey"`
	UserID    int `gorm:"primaryKey"`
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
	Name     string `json:"name" form:"name"`
	Avatar   string `json:"avatar" form:"avatar"`
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
	Name      string `json:"name" form:"name" binding:"required,lte=20" gorm:"size:255;index;unique;not null"`
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

type Image struct {
	ID         uint   `gorm:"primary_key" uri:"id" `
	CreatedOn  int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
	ModifiedOn int    `binding:"-" json:"modified_on,omitempty"`
	ArticleID  uint   `json:"-" form:"-" binding:"-"`
	Url        string `json:"url" form:"url" binding:"-"`
	ThumbUrl   string `json:"thumb_url" form:"thumb_url"`
	// Thumbnail int    `json:"-" form:"-"`
}

type SessionInfo struct {
	UserBase
	Unread   int        `json:"unread"`
	Messages []*Message `json:"messages"`
}
