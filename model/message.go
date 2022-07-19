package model

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type Message struct {
	ID        uint   `gorm:"primary_key" json:"id" `
	FromUid   int    `gorm:"index" json:"from_uid" form:"from_uid"`
	ToUid     int    `gorm:"index" json:"to_uid" form:"to_uid"`
	Content   string `json:"content" form:"content"`
	ImageUrl  string `json:"image_url" form:"image_url"`
	CreatedOn int    `gorm:"index" binding:"-" json:"created_on,omitempty"`
}
type MessageSession struct {
	ID         uint `gorm:"primary_key" json:"id" `
	Uid        int  `gorm:"index:uid_session,unique" json:"uid" form:"uid"`
	SessionId  int  `gorm:"index:uid_session,unique" json:"session_id" form:"session_id"`
	ModifiedOn int  `gorm:"index" binding:"-" json:"modified_on,omitempty"`
}

func SaveMessage(msg *Message) error {
	return db.Create(msg).Error
}

func SaveSession(session *MessageSession) error {
	return db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoUpdates: clause.AssignmentColumns([]string{"modified_on"})}).Create(session).Error
}

func GetSession(uid, pageNum, pageSize int) ([]*MessageSession, error) {
	var session []*MessageSession
	if err := db.Offset(pageNum).Limit(pageSize).Where("`uid` = ?", uid).Order("id desc").Find(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}
func GetMessages(fromUid, toUid, pageNum, pageSize int) ([]*Message, error) {
	var messages []*Message
	where := fmt.Sprintf("(`from_uid` = %d and `to_uid` = %d) or (`from_uid` = %d and `to_uid` = %d)", fromUid, toUid, toUid, fromUid)
	if err := db.Offset(pageNum).Limit(pageSize).Where(where).Order("id").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
