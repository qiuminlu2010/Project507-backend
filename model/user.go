package model

import (
	"gorm.io/gorm"
)

func GetUserTotal(maps interface{}) (int64, error) {
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserList(pageNum int, pageSize int, maps interface{}) ([]*User, error) {
	var users []*User
	err := db.Offset(pageNum).Where(maps).Limit(pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return users, nil
}

//TODO: GetUserArticle
func ExistUsername(username string) error {
	var user User
	return db.Where("username = ?", username).First(&user).Error
}

func ValidLogin(username string, password string) (User, error) {
	var user User
	if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func AddUser(user User) error {
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

func FollowUser(user User, followUser User) error {
	return db.Model(&user).Association("Follows").Append(&followUser)
}
