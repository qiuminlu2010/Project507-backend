package model

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"qiu/backend/pkg/setting"
)

var db *gorm.DB

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func Setup() {
	var (
		err         error
		tablePrefix string
	)

	sec := setting.DatabaseSetting
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	tablePrefix = sec.TablePrefix
	// log.Println("数据库初始化：", dbType, dbName, user, password, host, tablePrefix)
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	dsn1 := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sec.User, sec.HostMaster, sec.Name)
	dsn2 := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sec.User, sec.HostSlave1, sec.Name)
	dsn3 := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", sec.User, sec.HostSlave2, sec.Name)
	fmt.Println("master", sec.HostMaster, dsn1)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbName)
	db, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   tablePrefix,
		},
	})

	if err != nil {
		log.Println(err)
	}
	db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(dsn2), mysql.Open(dsn3)},
			Policy:   dbresolver.RandomPolicy{},
		}))
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
	db.Callback().Create().Before("gorm:create").Register("update_modify_time", updateTimeStampForUpdateCallback)
	db.Callback().Create().Before("gorm:create").Register("set_state", setState)
	db.Callback().Update().Before("gorm:update").Register("update_modify_time", updateTimeStampForUpdateCallback)

	err = db.SetupJoinTable(&Article{}, "LikedUsers", &ArticleLikeUsers{})
	if err != nil {
		panic(err)
	}
	err = db.SetupJoinTable(&User{}, "LikeArticles", &ArticleLikeUsers{})
	if err != nil {
		panic(err)
	}
	// db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		&User{},
		&Tag{},
		&Article{},
		&Image{},
		&Comment{},
		&Message{},
		&MessageSession{},
		// &Video{},
		// &Reply{},
	)
	if err != nil {
		log.Println("register table failed", err)
		os.Exit(0)
	}
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Tag{})
	// db.AutoMigrate(&Article{})
}
func setValue(db *gorm.DB, field string, value interface{}) {
	// field := "CreatedOn"
	if timeField := db.Statement.Schema.LookUpField(field); timeField != nil {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				if _, isZero := timeField.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); isZero {
					timeField.Set(db.Statement.Context, db.Statement.ReflectValue.Index(i), value)
				}
			}
		case reflect.Struct:
			if _, isZero := timeField.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
				timeField.Set(db.Statement.Context, db.Statement.ReflectValue, value)
			}
		}
	}
}
func setState(db *gorm.DB) {
	setValue(db, "State", 1)
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	setValue(db, "CreatedOn", time.Now().Unix())

}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	setValue(db, "ModifiedOn", time.Now().Unix())
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
