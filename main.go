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
	//鉴权接口
	http.HandleFunc("/addUser", utils.HandleInterceptor(api.AddUser))	//新增用户
	http.HandleFunc("/listUser", utils.HandleInterceptor(api.ListUser))	//列出用户
	http.HandleFunc("/deleteUser", utils.HandleInterceptor(api.DeleteUser))	//删除用户
	http.HandleFunc("/modifyUser", utils.HandleInterceptor(api.ModifyUser))	//修改用户
	http.HandleFunc("/getUser", utils.HandleInterceptor(api.GetUser))	//获取单个用户信息
	http.HandleFunc("/modifyPassword", utils.HandleInterceptor(api.ModifyPassword))	//修改密码

	//public接口
	http.HandleFunc("/join", utils.HandleInterceptor(api.Join))
	http.HandleFunc("/login", utils.HandleInterceptor(api.Login))

	//启动服务器
	log.Fatal(http.ListenAndServe(utils.HttpPort, nil))
}