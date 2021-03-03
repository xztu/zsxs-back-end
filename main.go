/*
   @Time : 2021/2/19 1:58 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : main
   @Description: 入口文件
*/

package main

import (
	"github.com/xztu/zsxs-back-end/commons/config"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/logger"
	"github.com/xztu/zsxs-back-end/commons/router"
)

// init 初始化操作
func init() {
	defer func() {
		// 捕获 PANIC
		if err := recover(); err != nil {
			logger.Log("捕获到了 PANIC 产生的异常. 未定义任何处理逻辑, 主进程已结束.")
			// 使用 defer 配合空 select 阻塞进程退出
			// 防止进程退出后, 容器被清理, 无法进行调试
			select {}
		}
	}()

	// 初始化数据库
	if err := orm.Init(); err != nil {
		// 初始化失败
		logger.Log("初始化数据库失败.")
		// 抛出异常
		logger.Panic(err)
	}

	// 打印初始化完成信息及版本号到控制台
	logger.Log("掌上忻师后端初始化完成, 当前版本 : " + config.Version)
}

// main 入口
func main() {
	// 初始化路由后启动监听
	if err := router.InitRouter().Run(); err != nil {
		// 处理错误
		logger.Panic(err)
	}
}
