package api

import (
	"fmt"
	"net/http"
	. "simple-blog/model"
	. "simple-blog/utils"
	"strconv"
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
func rightPassword(password string) bool {
	//合法的密码：
	//1. 至少8位字符
	//2. 必须包含字母数字
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
	return role >= int8(RoleUser) && role <= int8(RoleSuperAdmin)
}
//角色检查
func CheckRole(userRole int, r *http.Request) {
	cookie, _ := r.Cookie("login")
	id, _:= strconv.Atoi(cookie.Value)
	role, err:= GetRole(id)
	if err != nil {
		fmt.Printf("System Error: %v\n",err)
		panic(SYSTEM_ERROR)
	}
	if role == -1 {
		panic(ERROR_USERNAME_NOT_EXIST)
	}
	if role != userRole {
		panic(NO_POWER)
	}
}
//字段获取并检查
func MapToStruct(r *http.Request, checkType string, user *Account)  {
	params := RequestJsonInterface(r)
	//类型断言
	if username, ok := params["username"].(string); !ok {
		panic(ERROR_USERNAME_TYPE_WRONG)
	} else {
		if username == "" {
			panic(ERROR_USERNAME_NOT_NULL)
		}
		if !rightName(user.Username) {
			panic(ERROR_USERNAME_TYPE_WRONG)
		}
		status := CheckUserName(user.Username)
		if status != SUCCESS {
			panic(status)
		}
		user.Username = username
	}
	if nickname, ok := params["nickname"].(string); !ok {
		panic(ERROR_NICKNAME_TYPE_WRONG)
	} else {
		if nickname == "" {
			panic(ERROR_NICKNAME_NOT_NULL)
		}
		if !rightName(user.Nickname) {
			panic(ERROR_NICKNAME_TYPE_WRONG)
		}
		status := CheckNickName(user.Nickname)
		if status != SUCCESS {
			panic(status)
		}
		user.Nickname = nickname
	}
	if mail, ok := params["email"].(string); !ok {
		panic(ERROR_MAIL_TYPE_WRONG)
	} else {
		if mail == "" {
			panic(ERROR_MAIL_NOT_NULL)
		}
		status := CheckEmail(user.Email)
		if status != SUCCESS {
			panic(status)
		}
		user.Email = mail
	}

	//除了join，都需要role字段
	if checkType == "adduser" || checkType == "modifyuser" {
		if role, ok := params["role"].(float64); !ok {
			panic(ERROR_ROLE_TYPE_WRONG)
		} else {
			if !rightRole(user.Role) {
				panic(ERROR_ROLE_TYPE_WRONG)
			}
			user.Role = int8(role)
		}

	}
	//除了modifyuser，都需要password字段及创建时间
	if checkType == "adduser" || checkType == "join" {
		if password, ok := params["password"].(string); !ok {
			panic(ERROR_PASSWORD_TYPE_WRONG)
		} else {
			if password == "" {
				panic(ERROR_PASSWORD_NOT_NULL)
			}
			if !rightPassword(user.Password) {
				panic(ERROR_PASSWORD_TYPE_WRONG)
			}
			user.Password = password
		}
		user.CreateTime = time.Now().Unix()
	}
	//只有modifyuser需要id及更新时间
	if checkType == "modifyuser" {
		if id, ok := params["id"].(float64); !ok{
			panic(ERROR_USERID_TYPE_WRONG)
		} else {
			user.ID = int(id)
		}
		user.UpdateTime = time.Now().Unix()
	}
}
//错误处理
func errHandle(w http.ResponseWriter,err interface{})  {
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
		fmt.Printf("System Error: %v\n", result)
	}
}

// 添加用户
func AddUser(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		errHandle(w, err)
	}()
	CheckRole(RoleSuperAdmin, r)
	var user Account
	MapToStruct(r, "adduser", &user)
	status, err := CreateUser(&user)
	if err != nil {
		fmt.Printf("System Error: %v\n",err)
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
		errHandle(w, err)
	}()
	CheckRole(RoleSuperAdmin, r)
	params := RequestJsonInterface(r)
	var id int
	if tmp, ok := params["id"].(float64); !ok{
		panic(ERROR_USERID_TYPE_WRONG)
	} else {
		id = int(tmp)
	}
	if status := CheckUserID(id); status != SUCCESS {
		panic(status)
	}
	status , err := DeleteUserInDb(id)
	if err != nil {
		fmt.Printf("System Error: %v\n", err)
		panic(status)
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
		errHandle(w, err)
	}()
	var user Account
	MapToStruct(r, "modifyuser", &user)
	cookie, _ := r.Cookie("login")
	loginId, _ := strconv.Atoi(cookie.Value)
	status := CheckUserID(user.ID)
	if status != SUCCESS {
		panic(ERROR_USERNAME_NOT_EXIST)
	}
	if loginId != user.ID {
		CheckRole(RoleSuperAdmin, r)
	}
	status, err := EditUser(&user)
	if err != nil {
		fmt.Printf("System Error: %v\n", err)
		panic(status)
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
		errHandle(w, err)
	}()
	CheckRole(RoleSuperAdmin, r)
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
		fmt.Printf("System Error: %v\n", err)
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
		errHandle(w, err)
	}()
	CheckRole(RoleSuperAdmin, r)
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
		fmt.Printf("System Error: %v\n",err)
		panic(SYSTEM_ERROR)
	}
	w.Write(MapToBody(Map{
		"status":status,
		"desc": GetErrorMessage(status),
		"user":user,
	}))
}

//用户注册
func Join(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		errHandle(w, err)
	}()
	var user Account
	MapToStruct(r, "join", &user)
	user.Role = 1
	status, err := CreateUser(&user)
	if err != nil {
		fmt.Printf("System Error: %v\n",err)
		panic(status)
	}
	w.Write(MapToBody(Map{
		"status" : status,
		"desc" : GetErrorMessage(status),
		"id": user.ID,
	}))
}

//修改密码
func ModifyPassword(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		errHandle(w, err)
	}()
	var status, id int
	cookie, _ := r.Cookie("login")
	loginId, _ := strconv.Atoi(cookie.Value)
	var oldPassword, newPassword string
	params := RequestJsonInterface(r)
	if tmp, ok := params["id"].(float64); !ok{
		panic(ERROR_USERID_TYPE_WRONG)
	} else {
		id = int(tmp)
		if loginId != id {
			panic(NO_POWER)
		}
		status = CheckUserID(id)
		if status != SUCCESS {
			panic(ERROR_USERNAME_NOT_EXIST)
		}
	}
	if tmp, ok := params["oldpassword"].(string); !ok {
		panic(ERROR_PASSWORD_WRONG)
	} else {
		if tmp == "" {
			panic(ERROR_PASSWORD_WRONG)
		}
		oldPassword = tmp
	}
	if tmp, ok := params["newpassword"].(string); !ok {
		panic(ERROR_PASSWORD_TYPE_WRONG)
	} else {
		if tmp == "" {
			panic(ERROR_PASSWORD_TYPE_WRONG)
		}
		if !rightPassword(tmp) {
			panic(ERROR_PASSWORD_TYPE_WRONG)
		}
		newPassword = tmp
	}
	status, err := EditPassword(id, oldPassword, newPassword)
	if err != nil {
		fmt.Printf("System Error: %v\n",err)
		panic(status)
	}
	w.Write(MapToBody(Map{
		"status" : status,
		"desc" : GetErrorMessage(status),
	}))
}