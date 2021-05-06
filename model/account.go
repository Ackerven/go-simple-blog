package model

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
