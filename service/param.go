package service

type UserInfoParams struct {
	UserId   int    `uri:"user_id" json:"user_id" form:"user_id"`
	Username string `json:"username" form:"username"`
	Avator   string `json:"avator" form:"avator"`
}

type UserGetParams struct {
	UserId   int `json:"user_id" form:"user_id"`
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

// type UserIdParams
// type UserFollowsParams struct {
// 	UserInfoParams
// 	Follows []UserInfoParams `json:"follows"`
// }

type UpsertUserFollowParams struct {
	UserId   int `json:"user_id" form:"user_id"`
	FollowId int `json:"follow_id" form:"follow_id"`
	Type     int `json:"type" form:"type" binding:"gte=0,lte=1"`
}
type UserLoginParams struct {
	Id       uint   `uri:"id"`
	Username string `json:"username" form:"username" binding:"omitempty,printascii,gte=3,lte=20"`
	Password string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
	State    int    `json:"state" form:"state"`
}

type ArticleParams struct {
	Id         uint     `uri:"id"`
	UserID     uint     `json:"user_id" form:"user_id"`
	ImgName    []string `json:"-" form:"-" binding:"-"`
	TagName    []string `json:"tag_name" form:"tag_name"`
	TagID      []int    `json:"tag_id" form:"tag_id"`
	Title      string   `json:"title" form:"title"`
	Content    string   `json:"content" form:"content"`
	CreatedBy  string   `json:"created_by" form:"created_by"`
	ModifiedBy string   `json:"modified_by" form:"created_by"`
	State      int      `json:"state" form:"state" binding:"gte=0,lte=1"`
}

type ArticleLikeParams struct {
	ArticleId int `uri:"id"`
	UserId    int `json:"user_id" form:"user_id"`
	Type      int `json:"type" form:"type" binding:"gte=0,lte=1"`
}

type ArticleGetParams struct {
	Uid      int `json:"uid" form:"uid" binding:"gt=0"`
	PageNum  int `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int `json:"page_size" form:"page_size" binding:"gte=0"`
}

type ArticleAddParams struct {
	UserID  uint     `json:"user_id" form:"user_id" binding:"required,gt=0"`
	Title   string   `json:"title" form:"title" binding:"gte=0"`
	Content string   `json:"content" form:"content" binding:"gt=0"`
	TagName []string `json:"tag_name" form:"tag_name" `
	ImgName []string `json:"-" form:"-" binding:"-"`
}

// type ArticleGetParams struct {
// 	Uid      int `json:"uid" form:"uid"`
// 	PageNum  int `json:"page_num" form:"page_num"`
// 	PageSize int `json:"page_size" form:"page_size"`
// }

type CommentAddParams struct {
	UserId    int    `json:"user_id" form:"user_id" binding:"required,gt=0"`
	ArticleId int    `json:"article_id" form:"article_id" binding:"required,gt=0"`
	ReplyId   int    `json:"reply_id" form:"reply_id" binding:"gte=0"`
	Content   string `json:"content" form:"content" binding:"gt=0"`
}

type CommentGetParams struct {
	ArticleId int `json:"article_id" form:"article_id" binding:"required,gt=0"`
	PageNum   int `json:"page_num" form:"page_num"`
	PageSize  int `json:"page_size" form:"page_size"`
}

type ReplyAddParams struct {
	ReplyId int    `json:"comment_id" form:"comment_id" `
	UserId  int    `json:"user_id" form:"user_id" binding:"required,gt=0"`
	Content string `json:"content" form:"content" binding:"gt=0"`
}
