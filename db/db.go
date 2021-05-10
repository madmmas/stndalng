package db

import (
	"fmt"
	"stndalng/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func InitDB() {
	configuration := config.GetConfig()
	connect_string := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", configuration.DB_USERNAME, configuration.DB_PASSWORD, configuration.DB_HOST, configuration.DB_NAME)
	print(connect_string)
	db, err = gorm.Open("mysql", connect_string)
	if err != nil {
		print(err)
		panic("DB Connection Error")
	}
}

func DbManager() *gorm.DB {
	return db
}
