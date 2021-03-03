/*
   @Time : 2021/2/21 1:44 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : log
   @Description: 结构体 日志
*/

package structs

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	Source             string // 调用来源
	IP                 string // 用户 IP
	APIType            string // 接口类型
	OpenID             string // 用户 OpenID
	Method             string // 请求方法
	URL                string // 访问 URL
	Body               string // 访问 Body
	ResponseStatusCode int    // 响应状态码
}
