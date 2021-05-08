package api

import (
	"log"
	"net/http"
	. "simple-blog/model"
	. "simple-blog/utils"
	"time"
)


////检查用户名是否合法
//func rightName(username string) bool  {
//	str := []rune(username)
//	for _, v := range str {
//		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') || v == '_') {
//			return false
//		}
//	}
//	return true
//}

// 添加用户
func AddUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if code, ok := err.(int); ok {
				w.Write(MapToBody(Map{
					"status": code,
					"desc":   GetErrorMessage(code),
				}))
			} else {
				log.Fatal(err)
			}
		}
	}()
	var user Account
	var status int
	params := RequestJsonInterface(r)
	//类型断言
	if username, ok := params["username"].(string); !ok {
		panic(ERROR_USERNAME_TYPE_WRONG)
	} else {
		user.Username = username
	}
	if mail, ok := params["email"].(string); !ok {
		panic(ERROR_MAIL_TYPE_WRONG)
	} else {
		user.Email = mail
	}
	if nickname, ok := params["nickname"].(string); !ok {
		panic(ERROR_NICKNAME_TYPE_WRONG)
	} else {
		user.Nickname = nickname
	}
	if password, ok := params["password"].(string); !ok {
		panic(ERROR_PASSWORD_TYPE_WRONG)
	} else {
		user.Password = password
	}
	if role, ok := params["role"].(int8); !ok {
		panic(ERROR_ROLE_TYPE_WRONG)
	} else {
		user.Role = role
	}
	user.CreateTime = time.Now().Unix()

	////检查字段
	//if !rightName(user.Username) {
	//	panic(ERROR_USERNAME_TYPE_WRONG)
	//}
	//if !rightName(user.Nickname) {
	//	panic(ERROR_NICKNAME_TYPE_WRONG)
	//}


	status = CheckUserName(user.Username)
	if status != SUCCESS {
		panic(status)
	}
	status = CheckNickName(user.Nickname)
	if status != SUCCESS {
		panic(status)
	}
	status = CheckEmail(user.Email)
	if status != SUCCESS {
		panic(status)
	}


	status, err := CreateUser(&user)
	if err != nil {
		log.Fatal(err)
		panic(status)
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
