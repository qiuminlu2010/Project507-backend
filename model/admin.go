package model

import (
	"qiu/blog/pkg/e"

	"gorm.io/gorm"
)

func GetUserTotal(maps interface{}) (int64, error) {
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserList(pageNum int, pageSize int) ([]*User, error) {
	var users []*User
	err := db.Offset(pageNum).Limit(pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return users, nil
}

//TODO: GetUserArticle
func ExistUsername(username string) (res int64) {
	db.Where("username = ?", username).Count(&res)
	return
}

func ValidLogin(username string, password string) (*UserInfo, error) {

	var user UserInfo
	if err := db.Model(&User{}).Select("id", "username", "name", "avator").Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return nil, err
	}
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("follow_id = ?", user.ID).Count(&user.FanNum).Error; err != nil {
		return nil, err
	}
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("user_id = ?", user.ID).Count(&user.FollowNum).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AddUser(user *User) error {
	return db.Create(&user).Error
}

func UpdateUser(id uint, data interface{}) error {
	return db.Model(&User{}).Where("id = ?", id).Updates(data).Error

}

func DeleteUser(id uint) error {
	return db.Where("id = ?", id).Delete(&User{}).Error
}

func GetUsernameByID(id uint) string {
	var user User
	err := db.Select("username").Where("id = ?", id).First(&user).Error
	if err != nil {
		return ""
	}
	return user.Username
}
