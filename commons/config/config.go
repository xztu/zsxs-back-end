/*
   @Time : 2021/2/20 11:03 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : config
   @Description: 配置组件
*/

package config

import (
	"os"
	"strings"
)

// Version 版本号
// 用于 ORM 模块判断是否需要进行数据库结构自动迁移
// 用于输入版本信息到响应头
// 设置为距离当前时间较为久远的项目创始时间，避免每次调试时都需要自动同步数据库 可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var Version = "0000000 [ 2021/02/19 13:58:00 ]"

// IsDebugMode 是否运行调试模式
func IsDebugMode() bool {
	if strings.ToLower(os.Getenv("GIN_MODE")) == "release" {
		return false
	} else {
		return true
	}
}
