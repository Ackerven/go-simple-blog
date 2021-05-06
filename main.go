package main

import (
	"fmt"
	"log"
	"net/http"
	db "simple-blog/model"
	"simple-blog/utils"
)

func printf(w http.ResponseWriter, r *http.Request)  {
	var str string
	str = fmt.Sprintf("%s: %s", "mysql", utils.DbName)
	fmt.Fprintf(w, str)
}

func main() {
	//初始化数据库
	db.InitDb()

	//路由注册


	//启动服务器
	log.Fatal(http.ListenAndServe(utils.HttpPort, nil))
}