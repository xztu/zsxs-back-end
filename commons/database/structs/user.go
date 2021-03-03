/*
   @Time : 2021/2/21 6:13 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : user
   @Description: 结构体 用户
*/

package structs

import "gorm.io/gorm"

// User 用户
type User struct {
	gorm.Model
	Role     string `gorm:"not null"` // 用户角色
	OpenID   string `gorm:"not null"` // 微信账号 OpenID
	Username string `gorm:"not null"` // 用户名
	RealName string `gorm:"not null"` // 真实姓名
}

// UserLoginFailLog 用户 登陆失败日志
type UserLoginFailLog struct {
	gorm.Model
	Type     string `gorm:"not null"` // 错误类型
	OpenID   string `gorm:"not null"` // 微信账号 OpenID
	Username string `gorm:"not null"` // 用户名
	Password string `gorm:"not null"` // 密码
	SourceIp string `gorm:"not null"` // 访问 IP
}
