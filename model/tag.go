package model

//https://blog.csdn.net/weixin_45604257/article/details/105139862
type Tag struct {
	Model
	Name string `json:"name" form:"name" validate:"required,lte=20" `
	// Type       string `json:"type" form:"type" `
	CreatedBy  string `json:"created_by" form:"created_by" `
	ModifiedBy string `json:"modified_by" form:"modified_by"`
	State      int    `json:"state" form:"state"`
}

// func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
// 	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

// 	return
// }
func GetTags(pageNum int, pageSize int) (tags []Tag) {
	db.Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)

	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	return tag.ID > 0
}

func AddTag(tag Tag) error {
	return db.Create(&tag).Error
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	return tag.ID > 0
}

func GetTagIdByName(name string) (int, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil {
		return 0, err
	}
	return tag.ID, nil
}
func DeleteTag(id int) bool {
	return db.Where("id = ?", id).Delete(&Tag{}).RowsAffected > 0
}

func EditTag(id int, data interface{}) bool {
	return db.Model(&Tag{}).Where("id = ?", id).Updates(data).RowsAffected > 0
}
