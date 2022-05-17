package model

type User struct {
	Model

	Username string `json:"username"`
	Password string `json:"password"`
	State    int    `json:"state"`
}

func ExistUsername(username string) bool {
	var user User
	db.Select("id").Where("username = ?", username).First(&user)
	return user.ID > 0
}

func ValidLogin(username string, password string) bool {
	var user User
	db.Select("password").Where("username = ?", username).First(&user)

	return user.Password == password
}

func AddUser(username string, password string, state int) bool {
	db.Create(&User{
		Username: username,
		Password: password,
		State:    state,
	})

	return true
}

func ChangePassword(id int, data interface{}) bool {
	db.Model(&User{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteUser(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}
