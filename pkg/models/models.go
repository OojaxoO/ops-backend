package models 

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
  	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"

	"ops-backend/pkg/setting"
)

var DB *gorm.DB

func Setup () {
	var err error
	DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
    		    setting.DatabaseSetting.User,
    		    setting.DatabaseSetting.Password,
    		    setting.DatabaseSetting.Host,
				setting.DatabaseSetting.Name),
				), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Info),
					NamingStrategy: schema.NamingStrategy{
						// TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
						SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
					},
				})

	if err != nil {
      	fmt.Println(err)
      	return
  	}else {
  	    fmt.Println("connection succedssed")
	}
	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func CloseDB() {	

}