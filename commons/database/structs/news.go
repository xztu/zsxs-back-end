/*
   @Time : 2021/2/20 11:35 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : news
   @Description: 结构体 资讯
*/

package structs

import "gorm.io/gorm"

// NewsTopic 资讯 话题
type NewsTopic struct {
	gorm.Model
	Name  string `gorm:"not null"` // 名称
	Order uint   `gorm:"not null"` // 排序
	Icon  string `gorm:"not null"` // 图标
}

// NewsNews 资讯 资讯内容
type NewsNews struct {
	gorm.Model
	TopicID   uint   `gorm:"not null"` // 所属话题
	Order     uint   `gorm:"not null"` // 排序
	Headlines bool   `gorm:"not null"` // 头条
	Title     string `gorm:"not null"` // 标题
	Link      string `gorm:"not null"` // 链接
}
