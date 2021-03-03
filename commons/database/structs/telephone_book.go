/*
   @Time : 2021/2/24 10:09 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : telephone_book
   @Description: 结构体 电话簿
*/

package structs

import "gorm.io/gorm"

// TelephoneBookGroup 电话簿 分组
type TelephoneBookGroup struct {
	gorm.Model
	Type string // 类型
	Name string // 名称
	Rank uint   // 权重
}

// TelephoneBookContacts 电话簿 联系人
type TelephoneBookContacts struct {
	gorm.Model
	GroupID           uint   // 所属分组 ID
	Name              string // 姓名
	Office            string // 办公室电话
	OfficeShort       string // 办公室电话 小号
	ChinaMobile       string // 中国移动
	ChinaMobileShort  string // 中国移动 小号
	ChinaTelecom      string // 中国电信
	ChinaTelecomShort string // 中国电信 小号
	ChinaUnicom       string // 中国联通
	Remark            string // 备注
	Public            bool   // 是否公开
}
