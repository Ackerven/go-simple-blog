package model

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
	. "simple-blog/utils"
)

// 账户表
type Account struct {
	ID         int    `gorm:"primary_key" json:"id"`    // 主键
	Email      string `json:"email"`                    // 邮箱
	Nickname   string `json:"nickname" gorm:"not null"` // 昵称
	Username   string `json:"username" gorm:"not null"` // 用户名
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
		return ERROR_USERID_NOT_EXIET
	}
	return SUCCESS
}


//将用户写入数据库
func CreateUser(user *Account) (int,error) {
	var err error
	user.Password, err = ScryptPassword(user.Password)
	if err != nil {
		return ERROR, err
	}
	err = db.Create(&user).Error
	if err != nil {
		return ERROR, err
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
		return ERROR, err
	}
	return SUCCESS, nil
}
//编辑用户
func EditUser(id int, user *Account) (int,error) {
	var userMap = make(map[string]interface{})
	userMap["username"] = user.Username
	userMap["email"] = user.Email
	userMap["nickname"] = user.Nickname
	userMap["role"] = user.Role
	userMap["update_time"] = user.UpdateTime
	err = db.Model(&Account{}).Where("id = ?", id).Updates(userMap).Error
	if err != nil {
		return ERROR, err
	}
	return SUCCESS, nil
}
//获取单个用户的信息
func GetUserInDb(id int) (Account,error){
	var user Account
	err = db.Where("id = ?", id).First(&user).Error
	return user, err
}