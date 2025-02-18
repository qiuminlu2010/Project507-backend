package model

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type Message struct {
	ID        uint   `gorm:"primary_key" json:"id" `
	FromUid   int    `gorm:"index" json:"from_uid" form:"from_uid"`
	ToUid     int    `gorm:"index" json:"to_uid" form:"to_uid"`
	Content   string `gorm:"collate:utf8mb4" json:"content" form:"content" binding:"gte=1,lte=500"`
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

func GetAllSessions(uid int) ([]*MessageSession, error) {
	var session []*MessageSession
	if err := db.Where("`uid` = ?", uid).Order("id desc").Find(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func GetSession(uid, pageNum, pageSize int) ([]*MessageSession, error) {
	var session []*MessageSession
	if err := db.Offset(pageNum).Limit(pageSize).Where("`uid` = ?", uid).Order("id desc").Find(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func GetMessage(fromUid, toUid, pageNum, pageSize int) ([]*Message, error) {
	var messages []*Message
	where := fmt.Sprintf("(`from_uid` = %d and `to_uid` = %d) or (`from_uid` = %d and `to_uid` = %d)", fromUid, toUid, toUid, fromUid)
	if err := db.Offset(pageNum).Limit(pageSize).Where(where).Order("id desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// func GetMessages(fromUid int, toUids []int, pageNum int, pageSize int) ([]*Message, error) {
// 	var messages []*Message
// 	// where := fmt.Sprintf("(`from_uid` = %d and `to_uid` in (%v)) or (`from_uid` = %d and `to_uid` = %d)", fromUid, toUid, toUid, fromUid)
// 	if err := db.
// 	Offset(pageNum).
// 	Limit(pageSize).
// 	Where(db.Where("`from_uid` = ?", fromUid).Where("`to_uid` in (?)", toUids)).
// 	Or(db.Where("`to_uid` = ? ", fromUid).Where("`from_uid in (?)`", toUids)).
// 	Order("id desc").Find(&messages).Error; err != nil {
// 		return nil, err
// 	}
// 	return messages, nil
// }
