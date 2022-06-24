package service

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
	PageNum    int
	PageSize   int
}

type ArticleLikeParams struct {
	Id     int  `uri:"id"`
	UserID uint `json:"user_id" form:"user_id"`
}
