package model

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
	. "simple-blog/utils"
	"time"
)

// 账户表
type Account struct {
	ID         int    `gorm:"primary_key" json:"id"`    // 主键
	Username   string `json:"username" gorm:"not null"` // 用户名
	Nickname   string `json:"nickname" gorm:"not null"` // 昵称
	Email      string `json:"email"`                    // 邮箱
	Role       int8   `json:"role" gorm:"not null"`     // 角色
	Password   string // 密码
	CreateTime int64  `json:"create_time"` // 创建时间
	UpdateTime int64  `json:"update_time"` // 更新时间
}

//用户名是否存在
func CheckUserName(username string) int {
	var user Account
	db.Select("id").Where("username = ?", username).First(&user)
	if user.ID > 0 {
		return ERROR_USERNAME_EXIST
	}
	return SUCCESS
}
//邮箱是否存在
func CheckEmail(mail string) int {
	var user Account
	db.Select("id").Where("email = ?", mail).First(&user)
	if user.ID > 0 {
		return ERROR_MAIL_EXIST
	}
	return SUCCESS
}
//昵称是否存在
func CheckNickName(nickname string) int {
	var user Account
	db.Select("id").Where("nickname = ?", nickname).First(&user)
	if user.ID > 0 {
		return ERROR_NICKNAME_EXIST
	}
	return SUCCESS
}
//ID是否存在
func CheckUserID(id int) int  {
	var user Account
	err := db.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return ERROR_USERID_NOT_EXIST
	}
	return SUCCESS
}

//将用户写入数据库
func CreateUser(user *Account) (int,error) {
	var err error
	user.Password, err = ScryptPassword(user.Password)
	if err != nil {
		return SYSTEM_ERROR, err
	}
	err = db.Create(&user).Error
	if err != nil {
		return ERROR_DATABASE_WRITE, err
	}
	return SUCCESS, nil
}

//查询用户列表
func GetUserList(pageSize int, pageNum int) ([]Account, error) {
	var userList []Account
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&userList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return userList, nil
}

// ScryptPassword 密码加密 Scrypt算法
//scrypt https://pkg.go.dev/golang.org/x/crypto/scrypt
func ScryptPassword(password string) (string,error) {
	const KeyLen = 10
	//salt 盐
	salt := make([]byte, 8)
	salt = []byte{'I','L','O','V','E','Y','O','U'}

	HashPassword, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		return "", err
	}
	FinalPassword := base64.StdEncoding.EncodeToString(HashPassword)
	return FinalPassword, nil
}

//删除用户
func DeleteUserInDb(id int) (int,error) {
	err = db.Where("id = ?", id).Delete(&Account{}).Error
	if err != nil {
		return ERROR_DATABASE_DELETE, err
	}
	return SUCCESS, nil
}

//修改用户
func EditUser(user *Account) (int,error) {
	var userMap = make(map[string]interface{})
	userMap["username"] = user.Username
	userMap["email"] = user.Email
	userMap["nickname"] = user.Nickname
	userMap["role"] = user.Role
	userMap["update_time"] = user.UpdateTime
	err = db.Model(&Account{}).Where("id = ?", user.ID).Updates(userMap).Error
	if err != nil {
		return ERROR_DATABASE_WRITE, err
	}
	return SUCCESS, nil
}

//获取单个用户的信息
func GetUserInDb(id int) (Account,error){
	var user Account
	err = db.Where("id = ?", id).First(&user).Error
	return user, err
}

//修改密码
func EditPassword(id int, oldPassword, newPassword string) (int,error)  {
	user, err := GetUserInDb(id)
	if err != nil {
		return SYSTEM_ERROR, err
	}
	if tmp, _ := ScryptPassword(oldPassword); tmp != user.Password {
		return ERROR_PASSWORD_WRONG, nil
	}
	var userMap = make(map[string]interface{})
	tmp, err := ScryptPassword(newPassword)
	if err != nil {
		return SYSTEM_ERROR, err
	}
	userMap["password"] = tmp
	userMap["update_time"] = time.Now().Unix()
	err = db.Model(&Account{}).Where("id = ?", id).Updates(userMap).Error
	if err != nil {
		return SYSTEM_ERROR, err
	}
	return SUCCESS, nil
}

//登陆验证
func CheckLogin(username, password string, id *int) int {
	var user Account
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return ERROR_USERNAME_NOT_EXIST
	}
	*id = user.ID
	//fmt.Println(user.ID)
	if key, _ := ScryptPassword(password); key != user.Password {
		return ERROR_PASSWORD_WRONG
	}
	return SUCCESS
}

// 获取Role
func GetRole(id int) (int ,error){
	var user Account
	err = db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return -1, nil
	}
	return int(user.Role), nil
}