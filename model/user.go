package model

import (
	"qiu/blog/pkg/e"

	"gorm.io/gorm/clause"
)

func GetUsersByUsername(usernames []string) ([]*UserBase, error) {
	var users []*UserBase
	if err := db.Model(&User{}).Where("`username` in (?)", usernames).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByName(name string, pageNum int, pageSize int) ([]*UserBase, error) {
	var users []*UserBase
	if err := db.Model(&User{}).Where("`name` like ?", name).Offset(pageNum).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(userId uint) (*UserBase, error) {
	var user UserBase
	if err := db.Model(&User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func GetUserInfo(userId uint) (*UserInfo, error) {
	var user UserInfo
	if err := db.Model(&User{}).Select("id", "username", "name", "avatar").Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("follow_id = ?", userId).Count(&user.FanNum).Error; err != nil {
		return nil, err
	}
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("user_id = ?", userId).Count(&user.FollowNum).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func FollowUser(userId int, followId int) error {
	return db.Table(e.TABLE_USER_FOLLOWS).Clauses(clause.OnConflict{DoNothing: true}).Create(&UserIdFollowId{UserId: uint(userId), FollowId: uint(followId)}).Error
	// return db.Model(&user).Association("Follows").Append(&followUser)
}

func FollowUsers(userId uint, followIds []int) error {
	var group []UserIdFollowId
	for _, followId := range followIds {
		group = append(group, UserIdFollowId{UserId: userId, FollowId: uint(followId)})
	}
	return db.Table(e.TABLE_USER_FOLLOWS).Clauses(clause.OnConflict{DoNothing: true}).Create(group).Error
}

func UnFollowUser(userId int, followId int) error {
	return db.Table(e.TABLE_USER_FOLLOWS).Where("user_id = ?", userId).Where("follow_id = ?", followId).Delete(UserIdFollowId{}).Error
}

func UnFollowUsers(userId uint, followIds []int) error {
	return db.Table(e.TABLE_USER_FOLLOWS).Where("user_id = ?", userId).Where("follow_id in ?", followIds).Delete(UserIdFollowId{}).Error
}

func GetFollowIds(userId uint) ([]int, error) {
	// var followIds []*FollowId
	var followIds []int
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("`user_id` = ?", userId).Select("`follow_id`").Find(&followIds).Error; err != nil {
		return nil, err
	}
	return followIds, nil
}

func GetFollows(userId uint) ([]*User, error) {
	var user User
	user.ID = userId

	var follows []*User
	if err := db.Model(&user).Association("Follows").Find(&follows); err != nil {
		return nil, err
	}
	return follows, nil
}

func CountFollows(userId int) (int64, error) {
	var count int64
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserArticles(userId uint) ([]*UserArticlesCache, error) {
	var userArticlesCache []*UserArticlesCache
	if err := db.Model(&Article{}).Where("`owner_id` = ?", userId).Find(&userArticlesCache).Error; err != nil {
		return nil, err
	}
	return userArticlesCache, nil
}

func GetLikeArticleIds(userId uint) ([]int, error) {
	var likeArticleIds []int
	if err := db.Table(e.TABLE_ARTICLE_LIKE_USERS).Where("`user_id` = ?", userId).Select("`article_id`").Find(&likeArticleIds).Error; err != nil {
		return nil, err
	}
	return likeArticleIds, nil
}

func GetLikeArticles(userId uint) ([]*ArticleLikeUsers, error) {
	var likeArticles []*ArticleLikeUsers
	if err := db.Table(e.TABLE_ARTICLE_LIKE_USERS).Where("user_id = ?", userId).Find(&likeArticles).Error; err != nil {
		return nil, err
	}
	return likeArticles, nil
}

func GetFanIds(userId uint) ([]int, error) {
	// var followIds []*FollowId
	var fanIds []int
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("`follow_id` = ?", userId).Select("`user_id`").Find(&fanIds).Error; err != nil {
		return nil, err
	}
	return fanIds, nil
}

func GetFans(userId uint) ([]*UserBase, error) {
	var user User
	user.ID = userId

	var fans []*UserBase
	// if err := db.Model(&user).Association("Follows").Find(fans); err != nil {
	// 	return nil, err
	// }
	return fans, nil
}

func CountFans(userId int) (int64, error) {
	var count int64
	if err := db.Table(e.TABLE_USER_FOLLOWS).Where("follow_id = ?", userId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserId(username string) (uint, error) {
	var user User
	if err := db.Select("id").Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// func GetUserIds(usernames []string) ([]uint, error) {
// 	var users []User
// 	if err := db.Select("id").Where("username in (?)", usernames).Find(&users).Error; err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }
