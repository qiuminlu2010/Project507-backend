package model

type User struct {
	Model

	Username string `json:"username" form:"username" validate:"required,printascii,gte=6,lte=20"`
	Password string `json:"password" form:"password" validate:"required,printascii,gte=6,lte=20"`
	State    int    `json:"state" form:"state" validate:"gte=0,lte=1"`
}

func ExistUsername(username string) bool {
	var user User
	db.Select("id").Where("username = ?", username).First(&user)

	if user.ID > 0 {
		return true
	}
	return false

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

func ChangePassword(id int, data interface{}) bool {
	db.Model(&User{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteUser(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}
