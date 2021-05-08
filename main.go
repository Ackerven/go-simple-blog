package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"simple-blog/api"
	db "simple-blog/model"
	"simple-blog/utils"
)

func main() {
	//初始化数据库
	db.InitDb()

	//路由注册
	http.HandleFunc("/addUser", utils.HandleInterceptor(api.AddUser))

	//启动服务器
	log.Fatal(http.ListenAndServe(utils.HttpPort, nil))
}