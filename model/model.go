package model

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB

	username string = "root"
	password string = "xd_123456"
	dbName   string = "spiders"
)

func init() {
	var err error
	url := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbName)
	DB, err = gorm.Open("mysql", url)
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}
	log.Printf("open db ok[%s]\n", url)

	DB.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "sp_" + defaultTableName
	}

	tbname := "sp_douban_movie"
	log.Printf("has table %s [%t]\n", tbname, DB.HasTable(tbname))

}
