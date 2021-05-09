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
	http.HandleFunc("/listUser", utils.HandleInterceptor(api.ListUser))
	http.HandleFunc("/deleteUser", utils.HandleInterceptor(api.DeleteUser))
	http.HandleFunc("/modifyUser", utils.HandleInterceptor(api.ModifyUser))
	http.HandleFunc("/getUser", utils.HandleInterceptor(api.GetUser))


	//启动服务器
	log.Fatal(http.ListenAndServe(utils.HttpPort, nil))
}