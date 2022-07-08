package model

type ChatMsg struct {
	ID        uint   `gorm:"primary_key"  `
	FromUid   int    `json:"from_uid" form:"from_uid"`
	ToUid     int    `gorm:"index" json:"to_uid" form:"to_uid"`
	Content   string `json:"content" form:"content"`
	ImageUrl  string `json:"image_url" form:"image_url"`
	CreatedOn int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
}
