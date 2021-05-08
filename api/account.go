package api

import (
	"fmt"
	"log"
	"net/http"
	. "simple-blog/model"
	. "simple-blog/utils"
	"time"
)


//检查用户名及昵称是否合法
func rightName(username string) bool  {
	str := []rune(username)
	for _, v := range str {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') || v == '_') {
			return false
		}
	}
	return true
}
//检查密码是否合法
//合法的密码：
//1. 至少8位字符
//2. 必须包含字母数字
func rightPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	str := []rune(password)
	countNumber := 0
	countChar := 0
	for _, v := range str {
		if v < '!' || v > '~' {
			return false
		}
		if v >= '0' && v <= '9' {
			countNumber++
		}
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			countChar++
		}
	}
	if countNumber == 0 || countChar == 0 {
		return false
	}
	return true
}
//用户类型检查
func rightRole(role int8) bool {
	return role > 0 && role < 4
}

// 添加用户
func AddUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var user Account
	var status int
	params := RequestJsonInterface(r)
	//类型断言
	if username, ok := params["username"].(string); !ok {
		HandleError(ERROR_USERNAME_TYPE_WRONG, w, r)
		_ = params["username"].(string)
		return
	} else {
		if username == "" {
			HandleError(ERROR_USERNAME_NOT_NULL, w, r)
			return
		}
		user.Username = username
	}
	if mail, ok := params["email"].(string); !ok {
		HandleError(ERROR_MAIL_TYPE_WRONG, w, r)
		_ = params["email"].(string)
		return
	} else {
		if mail == "" {
			HandleError(ERROR_MAIL_NOT_NULL, w, r)
			return
		}
		user.Email = mail
	}
	if nickname, ok := params["nickname"].(string); !ok {
		HandleError(ERROR_NICKNAME_TYPE_WRONG, w, r)
		_ = params["nickname"].(string)
		return
	} else {
		if nickname == "" {
			HandleError(ERROR_NICKNAME_NOT_NULL, w, r)
			return
		}
		user.Nickname = nickname
	}
	if password, ok := params["password"].(string); !ok {
		HandleError(ERROR_PASSWORD_TYPE_WRONG, w, r)
		_ = params["password"].(string)
		return
	} else {
		if password == "" {
			HandleError(ERROR_PASSWORD_NOT_NULL, w, r)
			return
		}
		user.Password = password
	}
	if role, ok := params["role"].(float64); !ok {
		HandleError(ERROR_ROLE_TYPE_WRONG, w, r)
		_ = params["role"].(int8)
		return
	} else {
		user.Role = int8(role)
	}
	user.CreateTime = time.Now().Unix()

	//检查字段
	if !rightName(user.Username) {
		HandleError(ERROR_USERNAME_TYPE_WRONG, w, r)
		return
	}
	if !rightName(user.Nickname) {
		HandleError(ERROR_NICKNAME_TYPE_WRONG, w, r)
		return
	}
	if !rightPassword(user.Password) {
		HandleError(ERROR_PASSWORD_TYPE_WRONG, w, r)
		return
	}
	if !rightRole(user.Role) {
		HandleError(ERROR_ROLE_TYPE_WRONG, w, r)
		return
	}

	//数据库是否有记录
	status = CheckUserName(user.Username)
	if status != SUCCESS {
		HandleError(status, w, r)
		return
	}
	status = CheckNickName(user.Nickname)
	if status != SUCCESS {
		HandleError(status, w, r)
		return
	}
	status = CheckEmail(user.Email)
	if status != SUCCESS {
		HandleError(status, w, r)
		return
	}

	status, err := CreateUser(&user)
	if err != nil {
		HandleError(status, w, r)
		log.Fatal(err)
	}
	w.Write(MapToBody(Map{
		"status" : status,
		"desc" : GetErrorMessage(status),
		"id": user.ID,
	}))
}

// 删除用户
func DeleteUser() {

}

// 修改用户
func ModifyUser() {

}

// 查询用户
func GetUser()  {

}

// 列出用户
func ListUser()  {

}
