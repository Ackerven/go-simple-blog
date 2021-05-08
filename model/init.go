package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"simple-blog/utils"
	"time"
)

var (
	db *gorm.DB
	err error
)

func InitDb() {
Conn:
	db, err = gorm.Open(utils.Db, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	))
	if err != nil {
		fmt.Println("[!] 数据库连接失败！尝试重新连接...", err)
		time.Sleep(time.Second * 5)
		goto Conn
	}

	db.DB().SetMaxIdleConns(30)
	db.DB().SetMaxOpenConns(500)
	//db.DB().SetConnMaxLifetime(time.Second * 10)
	// 取消自动加复数
	db.SingularTable(true)
	db.AutoMigrate(&Account{}, &Blog{})
	//db.Close()
}