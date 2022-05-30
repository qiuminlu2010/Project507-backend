package model

type User struct {
	Model
	Username  string    `json:"username" form:"username" binding:"omitempty,printascii,gte=6,lte=20" gorm:"unique"`
	Password  string    `json:"password" form:"password" binding:"omitempty,printascii,gte=6,lte=20"`
	StudentId string    `json:"student_id" form:"student_id" binding:"omitempty,numeric"`
	Articles  []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	State     int       `json:"state" form:"state" binding:"gte=0,lte=1"`
}

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
