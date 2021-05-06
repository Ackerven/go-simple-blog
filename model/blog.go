package model

// 博客表
type Blog struct {
	ID         int    `gorm:"primary_key" json:"id"`             // 主键
	AuthorID   int    `json:"author_id"`                         // 作者id
	Title      string `json:"title" gorm:"not null"`             // 标题
	Content    string `json:"content" gorm:"type:text;not null"` // 正文
	Draft      bool   `json:"draft" gorm:"default:true"`         // 是否草稿
	Private    bool   `json:"private" gorm:"default:false"`      // 是否私有
	CreateTime int64  `json:"create_time"`                       // 创建时间
	UpdateTime int64  `json:"update_time"`                       // 更新时间
}

