package model

type User struct {
	Model

	Username  string `json:"username" form:"username" validate:"omitempty,printascii,gte=6,lte=20" gorm:"unique"`
	Password  string `json:"password" form:"password" validate:"omitempty,printascii,gte=6,lte=20"`
	StudentId string `json:"student_id" form:"student_id" validate:"omitempty,numeric"`
	State     int    `json:"state" form:"state" validate:"gte=0,lte=1"`
}

func ExistUsername(username string) bool {
	var user User
	db.Select("id").Where("username = ?", username).First(&user)

	return user.ID > 0

}

func ValidLogin(username string, password string) (bool, error) {
	var user User
	err := db.Select("password").Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, err
	}
	if user.Password != password {
		return false, nil
	}
	return true, nil
}

func AddUser(user User) error {
	err := db.Create(&user).Error

	return err
}

func UpdatePassword(id uint, data interface{}) bool {
	return db.Model(&User{}).Where("id = ?", id).Updates(data).RowsAffected > 0

}

func DeleteUser(id uint) bool {
	return db.Where("id = ?", id).Delete(&User{}).RowsAffected > 0
}

func GetUsernameByID(id uint) string {
	var user User
	err := db.Select("username").Where("id = ?", id).First(&user).Error
	if err != nil {
		return ""
	}
	return user.Username
}
