package api_param

type PageGetParams struct {
	PageNum  int `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int `json:"page_size" form:"page_size" binding:"gte=0"`
}

type UserInfoParams struct {
	UserId   int    `uri:"user_id" json:"user_id" form:"user_id"`
	Username string `json:"username" form:"username"`
	Avatar   string `json:"avatar" form:"avatar"`
}

type UserListGetParams struct {
	Name     string `json:"name" form:"name" binding:"gt=0"`
	PageNum  int    `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int    `json:"page_size" form:"page_size" binding:"gte=0"`
}

type UserLoginParams struct {
	Username string `json:"username" form:"username" binding:"omitempty,printascii,gte=4,lte=20"`
	Password string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
}

type UsersGetParams struct {
	Name     string `json:"name" form:"name" binding:"gt=0"`
	PageNum  int    `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int    `json:"page_size" form:"page_size" binding:"gte=0"`
}

type UserAddParams struct {
	Username string `json:"username" form:"username" binding:"printascii,gte=4,lte=20"`
	Name     string `json:"name" form:"name" binding:"gte=0"`
	Password string `json:"password" form:"password" binding:"printascii,gte=6,lte=100"`
}

type UserUpdateParams struct {
	// Username string `json:"username" form:"username" binding:"omitempty,printascii,gte=4,lte=20"`
	UserId   int    `json:"user_id" form:"user_id" `
	Name     string `json:"name" form:"name" binding:"omitempty,gte=4,lte=20"`
	Password string `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=100"`
	State    int    `json:"state" form:"state" binding:"omitempty,gte=0,lte=1"`
}

type FollowsGetParams struct {
	UserId   int `json:"user_id" form:"user_id"`
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

type UpsertUserFollowParams struct {
	UserId   int `json:"user_id" form:"user_id"`
	FollowId int `json:"follow_id" form:"follow_id"`
	Type     int `json:"type" form:"type" binding:"gte=0,lte=1"`
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
	Uid      int `json:"uid" form:"uid" binding:"gte=0"`
	PageNum  int `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int `json:"page_size" form:"page_size" binding:"gte=0"`
	Type     int `json:"type" form:"type"`
}

type ArticleAddParams struct {
	UserID  uint     `json:"user_id" form:"user_id" binding:"required,gt=0"`
	Title   string   `json:"title" form:"title" binding:"gte=0"`
	Content string   `json:"content" form:"content" binding:"gt=0"`
	TagName []string `json:"tag_name" form:"tag_name" `
	ImgUrl  []string `json:"-" form:"-" binding:"-"`
}

type ArticleUpdateParams struct {
	ArticleId int
	Title     string `json:"title" form:"title" binding:"omitempty,gt=0"`
	Content   string `json:"content" form:"content" binding:"omitempty,gt=0"`
}

type ArticleAddTagsParams struct {
	ArticleId int
	TagName   []string `json:"tag_name" form:"tag_name" `
}

type TagsGetParams struct {
	TagName  string `json:"tag_name" form:"tag_name" binding:"gt=0"`
	PageNum  int    `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int    `json:"page_size" form:"page_size" binding:"gte=0"`
}

type TagArticleGetParams struct {
	Uid      int    `json:"uid" form:"uid" binding:"gte=0"`
	TagName  string `json:"tag_name" form:"tag_name" binding:"gte=0"`
	PageNum  int    `json:"page_num" form:"page_num" binding:"gte=0"`
	PageSize int    `json:"page_size" form:"page_size" binding:"gte=0"`
}

type CommentAddParams struct {
	UserId    int    `json:"user_id" form:"user_id" binding:"required,gt=0"`
	ArticleId int    `json:"article_id" form:"article_id" binding:"required,gt=0"`
	ReplyId   int    `json:"reply_id" form:"reply_id" binding:"gte=0"`
	Content   string `json:"content" form:"content" binding:"gt=0"`
}

type CommentGetParams struct {
	ArticleId int `json:"article_id" form:"article_id" binding:"required,gt=0"`
	UserId    int `json:"user_id" form:"user_id" `
	PageNum   int `json:"page_num" form:"page_num"`
	PageSize  int `json:"page_size" form:"page_size"`
}

type LikeCommentParams struct {
	UserId    int `json:"user_id" form:"user_id" binding:"required,gt=0"`
	CommentId int `json:"comment_id" form:"comment_id" `
	Type      int `json:"type" form:"type" binding:"gte=0,lte=1"`
}

type RefreshTokenParams struct {
	UserId int    `json:"user_id" form:"user_id" binding:"omitempty,gt=0"`
	Uuid   string `json:"uuid" form:"uuid"`
}

// type MsgClientParams struct {
// 	FromUid int `json:"from_uid" form:"from_uid" uri:"from_uid"`
// 	// ToUid   int `json:"to_uid" form:"to_uid" uri:"to_uid"`
// }

type MessageGetParams struct {
	Offset  int `json:"offset" form:"offset"`
	Limit   int `json:"limit" form:"limit"`
	FromUid int `json:"from_uid" form:"from_uid" uri:"from_uid"`
	ToUid   int `json:"to_uid" form:"to_uid" uri:"to_uid"`
}

type SessionGetParams struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
	Uid      int `json:"uid" form:"uid" uri:"uid"`
}

type UpdateUnReadMessageParams struct {
	Uid       int `json:"uid" form:"uid" `
	SessionId int `json:"session_uid" form:"session_uid"`
}
