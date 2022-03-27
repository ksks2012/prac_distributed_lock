package model

import (
	"fmt"
	"time"

	"github.com/distributed_lock/global"
	"github.com/distributed_lock/pkg/setting"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
// TODO: rewrite
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// process public field
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	addr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: databaseSetting.TablePrefix,   // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		return nil, err
	}

	if global.AppSetting.RunMode == "debug" {
		db.Debug()
	}

	sqlDB, err_conn := db.DB()
	if err_conn != nil {
		fmt.Println("get db failed:", err)
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Duration(databaseSetting.MaxLifetime) * time.Second)
	// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	// db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
