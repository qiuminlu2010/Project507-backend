package model

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"qiu/blog/pkg/setting"
)

var db *gorm.DB

type Model struct {
	ID         uint           `gorm:"primary_key" uri:"id" `
	CreatedOn  time.Time      `binding:"-" json:"-"`
	ModifiedOn time.Time      `binding:"-" json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index"  binding:"-" json:"-"`
}

func Setup() {
	var (
		err                                       error
		dbName, user, password, host, tablePrefix string
	)

	sec := setting.DatabaseSetting
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbName = sec.Name
	user = sec.User
	password = sec.Password
	host = sec.Host
	tablePrefix = sec.TablePrefix
	// log.Println("数据库初始化：", dbType, dbName, user, password, host, tablePrefix)
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   tablePrefix,
		},
	})
	// db, err = gorm.Open(dbType, dsn)

	if err != nil {
		log.Println(err)
	}

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return tablePrefix + defaultTableName
	// }
	sqlDB, err := db.DB()
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	// db.SingularTable(true)
	// db.LogMode(true)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	// db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.Callback().Create().Before("gorm:create").Register("update_create_time", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("update_modify_time", updateTimeStampForUpdateCallback)
	err = db.AutoMigrate(
		&User{},
		&Tag{},
		&Article{},
		&Image{},
	)
	if err != nil {
		log.Println("register table failed", err)
		os.Exit(0)
	}
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Tag{})
	// db.AutoMigrate(&Article{})
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	field := "CreatedOn"
	if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				if _, isZero := timeField.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); isZero {
					timeField.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), time.Now())
				}
			}
		case reflect.Struct:
			if _, isZero := timeField.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
				timeField.Set(db.Statement.Context, db.Statement.ReflectValue, time.Now())
			}
		}
	}

}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	field := "ModifiedOn"
	if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				timeField.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), time.Now())
			}
		case reflect.Struct:
			timeField.Set(db.Statement.Context, db.Statement.ReflectValue, time.Now())
		}
	}

}

// func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
// 	if _, ok := scope.Get("gorm:update_column"); !ok {
// 		scope.SetColumn("ModifiedOn", time.Now())
// 	}
// }

// func deleteCallback(scope *gorm.Scope) {
// 	if !scope.HasError() {
// 		var extraOption string
// 		if str, ok := scope.Get("gorm:delete_option"); ok {
// 			extraOption = fmt.Sprint(str)
// 		}

// 		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

// 		if !scope.Search.Unscoped && hasDeletedOnField {
// 			scope.Raw(fmt.Sprintf(
// 				"UPDATE %v SET %v=%v%v%v",
// 				scope.QuotedTableName(),
// 				scope.Quote(deletedOnField.DBName),
// 				scope.AddToVars(time.Now()),
// 				addExtraSpaceIfExist(scope.CombinedConditionSql()),
// 				addExtraSpaceIfExist(extraOption),
// 			)).Exec()
// 		} else {
// 			scope.Raw(fmt.Sprintf(
// 				"DELETE FROM %v%v%v",
// 				scope.QuotedTableName(),
// 				addExtraSpaceIfExist(scope.CombinedConditionSql()),
// 				addExtraSpaceIfExist(extraOption),
// 			)).Exec()
// 		}
// 	}
// }

// func addExtraSpaceIfExist(str string) string {
// 	if str != "" {
// 		return " " + str
// 	}
// 	return ""
// }

// func CloseDB() {
// 	defer db.Close()
// }
