package api

import (
	"fmt"
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
		err := recover()
		if err == nil {
			return
		}
		switch result := err.(type) {
		case int :
			//status := err.(int)
			w.Write(MapToBody(Map{
				"status":result,
				"desc": GetErrorMessage(result),
			}))
		default:
			fmt.Printf("系统错误：%v", result)
		}
	}()
	var user Account
	var status int
	params := RequestJsonInterface(r)
	//类型断言
	if username, ok := params["username"].(string); !ok {
		panic(ERROR_USERNAME_TYPE_WRONG)
	} else {
		if username == "" {
			panic(ERROR_USERNAME_NOT_NULL)
		}
		user.Username = username
	}
	if mail, ok := params["email"].(string); !ok {
		panic(ERROR_MAIL_TYPE_WRONG)
	} else {
		if mail == "" {
			panic(ERROR_MAIL_NOT_NULL)
		}
		user.Email = mail
	}
	if nickname, ok := params["nickname"].(string); !ok {
		panic(ERROR_NICKNAME_TYPE_WRONG)
	} else {
		if nickname == "" {
			panic(ERROR_NICKNAME_NOT_NULL)
		}
		user.Nickname = nickname
	}
	if password, ok := params["password"].(string); !ok {
		panic(ERROR_PASSWORD_TYPE_WRONG)
	} else {
		if password == "" {
			panic(ERROR_PASSWORD_NOT_NULL)
		}
		user.Password = password
	}
	if role, ok := params["role"].(float64); !ok {
		panic(ERROR_ROLE_TYPE_WRONG)
	} else {
		user.Role = int8(role)
	}
	user.CreateTime = time.Now().Unix()

	//检查字段
	if !rightName(user.Username) {
		panic(ERROR_USERNAME_TYPE_WRONG)
	}
	if !rightName(user.Nickname) {
		panic(ERROR_NICKNAME_TYPE_WRONG)
	}
	if !rightPassword(user.Password) {
		panic(ERROR_PASSWORD_TYPE_WRONG)
	}
	if !rightRole(user.Role) {
		panic(ERROR_ROLE_TYPE_WRONG)
	}

	//数据库是否有记录
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
		fmt.Sprintf("系统错误：%v",err)
		panic(status)
	}
	w.Write(MapToBody(Map{
		"status" : status,
		"desc" : GetErrorMessage(status),
		"id": user.ID,
	}))
}

// 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		switch result := err.(type) {
		case int :
			//status := err.(int)
			w.Write(MapToBody(Map{
				"status":result,
				"desc": GetErrorMessage(result),
			}))
		default:
			fmt.Printf("系统错误：%v", result)
		}
	}()
	params := RequestJsonInterface(r)
	var id int
	if tmp, ok := params["id"].(float64); !ok{
		panic(ERROR_USERID_TYPE_WRONG)
	} else {
		id = int(tmp)
	}
	status , err := DeleteUserInDb(id)
	if err != nil {
		fmt.Printf("系统错误：%v", err)
		panic(ERROR_DATABASE_DELETE)
	}
	w.Write(MapToBody(Map{
		"status": status,
		"desc": GetErrorMessage(status),
	}))
}

// 修改用户
func ModifyUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		switch result := err.(type) {
		case int :
			//status := err.(int)
			w.Write(MapToBody(Map{
				"status":result,
				"desc": GetErrorMessage(result),
			}))
		default:
			fmt.Printf("系统错误：%v", result)
		}
	}()
	var user Account
	var status, id int
	params := RequestJsonInterface(r)

	//类型断言
	if username, ok := params["username"].(string); !ok {
		panic(ERROR_USERNAME_TYPE_WRONG)
	} else {
		if username == "" {
			panic(ERROR_USERNAME_NOT_NULL)
		}
		user.Username = username
	}
	if mail, ok := params["email"].(string); !ok {
		panic(ERROR_MAIL_TYPE_WRONG)
	} else {
		if mail == "" {
			panic(ERROR_MAIL_NOT_NULL)
		}
		user.Email = mail
	}
	if nickname, ok := params["nickname"].(string); !ok {
		panic(ERROR_NICKNAME_TYPE_WRONG)
	} else {
		if nickname == "" {
			panic(ERROR_NICKNAME_NOT_NULL)
		}
		user.Nickname = nickname
	}
	if role, ok := params["role"].(float64); !ok {
		panic(ERROR_ROLE_TYPE_WRONG)
	} else {
		user.Role = int8(role)
	}

	//数据库是否有记录
	status = CheckUserName(user.Username)
	if status != SUCCESS {
		panic(status)
		return
	}
	status = CheckNickName(user.Nickname)
	if status != SUCCESS {
		panic(status)
		return
	}
	status = CheckEmail(user.Email)
	if status != SUCCESS {
		panic(status)
		return
	}

	if tmp, ok := params["id"].(float64); !ok{
		panic(ERROR_USERID_TYPE_WRONG)
		return
	} else {
		id = int(tmp)
	}
	status, err := EditUser(id, &user)
	if err != nil {
		fmt.Printf("系统错误：%v", err)
		panic(ERROR_DATABASE_WRITE)
	}
	w.Write(MapToBody(Map{
		"status": status,
		"desc": GetErrorMessage(status),
	}))
}


// 列出用户
func ListUser(w http.ResponseWriter, r *http.Request)  {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		switch result := err.(type) {
		case int :
			//status := err.(int)
			w.Write(MapToBody(Map{
				"status":result,
				"desc": GetErrorMessage(result),
			}))
		default:
			fmt.Printf("系统错误：%v", result)
		}
	}()
	params := RequestJsonInterface(r)
	var pageSize, pageNum int
	if tmp, ok := params["pagesize"].(float64); !ok {
		panic(ERROR_PAGESIZE_TYPE_WRONG)
	} else {
		pageSize = int(tmp)
	}
	if tmp, ok := params["pagenum"].(float64); !ok {
		panic(ERROR_PAGENUM_TYPE_WRONG)
	} else {
		pageNum = int(tmp)
	}

	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}
	userList, err := GetUserList(pageSize, pageNum)
	if err != nil {
		fmt.Printf("系统错误：%v", err)
		panic(ERROR_DATABASE_SEARCH)
	}
	w.Write(MapToBody(Map{
		"status":SUCCESS,
		"desc": GetErrorMessage(SUCCESS),
		"result":userList,
	}))
}

// 查询用户
func GetUser(w http.ResponseWriter, r *http.Request)  {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		switch result := err.(type) {
		case int :
			//status := err.(int)
			w.Write(MapToBody(Map{
				"status":result,
				"desc": GetErrorMessage(result),
			}))
		default:
			fmt.Printf("系统错误：%v", result)
		}
	}()
	params := RequestJsonInterface(r)
	var id, status int
	if tmp, ok := params["id"].(float64); !ok{
		panic(ERROR_USERID_TYPE_WRONG)
	} else {
		id = int(tmp)
	}

	if status = CheckUserID(id); status != SUCCESS {
		panic(status)
	}

	user, err := GetUserInDb(id)
	if err != nil {
		fmt.Sprintf("系统错误：%v",err)
		panic(ERROR)
	}
	w.Write(MapToBody(Map{
		"status":status,
		"desc": GetErrorMessage(status),
		"user":user,
	}))
}
