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
	Model
	ArticleID uint   `json:"-" form:"-" binding:"-"`
	Filename  string `json:"filename" form:"filename" binding:"-"`
}
type Article struct {
	Model
	OwnerID   uint    `gorm:"index" json:"owner_id"`
	User      User    `gorm:"foreignkey:OwnerID" binding:"-" json:"-"` // 使用 OwnerID  作为外键
	Tags      []Tag   `gorm:"many2many:article_tags;" json:"tags"`
	Images    []Image `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignkey:ArticleID" json:"images"`
	Title     string  `json:"title" form:"title"`
	Content   string  `json:"content" form:"content"`
	LikeCount int64   `json:"like_count" form:"like_count" binding:"-"`
	// Collect int     `json:"collect" form:"collect" binding:"-"`
	// Watch   int     `json:"watch" form:"watch" binding:"-"`
	LikedUsers []User `gorm:"many2many:article_like_users;" json:"-"`
	IsLike     bool   `json:"is_like" form:"is_like" binding:"-"`
	//TODO: Comments   []Comment
	// CreatedBy  string `json:"-" form:"created_by" binding:"-"`
	// ModifiedBy string `json:"-" form:"created_by" binding:"-"`
	State int `json:"state" form:"state" binding:"-"`
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
	UserID  uint `json:"user_id"`
	Content uint `json:"content" form:"content"`
}

type User struct {
	Model
	Username  string `json:"username" form:"username" binding:"omitempty,printascii,gte=6,lte=20" gorm:"unique"`
	Password  string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
	StudentId string `json:"student_id" form:"student_id" binding:"omitempty,numeric"`
	// Articles     []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"articles,omitempty" binding:"-"`
	LikeArticles []*Article `gorm:"many2many:article_like_users" binding:"-" json:"like_articles"`
	Follows      []*User    `gorm:"many2many:user_follows"`
	// FollowNum    int        `json:"follow_num" form:"follow_num" binding:"-"`
	// FanNum       int        `json:"fan_num" form:"fan_num" binding:"-"`
	Avator string `json:"avator" form:"avator"`
	State  int    `json:"state" form:"state" binding:"gte=0,lte=1"`
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

type ArticleIdUserId struct {
	ArticleId uint `json:"article_id"`
	UserId    uint `json:"user_id"`
}

type UserIdFollowId struct {
	UserId   uint `json:"user_id"`
	FollowId uint `json:"follow_id"`
}
type Tag struct {
	Model
	Name string `json:"name" form:"name" binding:"required,lte=20" gorm:"unique;not null"`
	//Type       string `json:"type" form:"type" `
	CreatedBy  string    `json:"-" form:"created_by" binding:"-" `
	ModifiedBy string    `json:"-" form:"modified_by" binding:"-" `
	State      int       `json:"-" form:"state" binding:"-" `
	Articles   []Article `gorm:"many2many:article_tags;" binding:"-" json:"-"`
}
