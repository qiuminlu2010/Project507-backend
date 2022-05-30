package model

//https://blog.csdn.net/weixin_45604257/article/details/105139862
type Tag struct {
	Model
	Name string `json:"name" form:"name" binding:"required,lte=20" gorm:"unique;not null"`
	//Type       string `json:"type" form:"type" `
	CreatedBy  string    `json:"-" form:"created_by" binding:"-" `
	ModifiedBy string    `json:"-" form:"modified_by" binding:"-" `
	State      int       `json:"-" form:"state" binding:"-" `
	Articles   []Article `gorm:"many2many:article_tags;" binding:"-" json:"-"`
}

// func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
// 	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

// 	return
// }
func GetTags(pageNum int, pageSize int) (tags []Tag) {
	db.Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int64) {
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

func GetTagIdByName(name string) (uint, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil {
		return 0, err
	}
	return tag.ID, nil
}
func DeleteTag(id uint) error {
	var tag Tag
	tag.ID = id

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Model(&tag).Association("Articles").Clear(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id = ?", id).Delete(Tag{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func RecoverTag(id uint) error {
	var tag Tag
	tag.ID = id
	if err := db.Unscoped().Model(&tag).Update("deleted_at", nil).Error; err != nil {
		return err
	}
	return nil
}

func EditTag(id uint, data interface{}) error {
	return db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
}
