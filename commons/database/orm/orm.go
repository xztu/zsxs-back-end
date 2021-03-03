/*
   @Time : 2021/2/20 11:06 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : orm
   @Description: 对象关系映射组件
*/

package orm

import (
	"errors"
	"fmt"
	"github.com/cengsin/oracle"
	"github.com/xztu/zsxs-back-end/commons/config"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/commons/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

// 对外暴露的数据库实例
var MySQL *gorm.DB
var Oracle *gorm.DB

// Init 初始化数据库
func Init() error {
	// 校验是否配置 MYSQL_DSN
	if os.Getenv("MYSQL_DSN") == "" {
		// 未配置，返回错误
		return errors.New("未配置 MYSQL_DSN , 请在环境变量中配置 MYSQL_DSN ( 例 : user:password@tcp(hostname)/database?charset=utf8mb4&parseTime=True&loc=Local )")
	}
	// 校验是否配置 ORACLE_DSN
	if os.Getenv("ORACLE_DSN") == "" {
		// 未配置，返回错误
		return errors.New("未配置 ORACLE_DSN , 请在环境变量中配置 ORACLE_DSN ( 例 : user/password@hostname:port/service_name )")
	}
	// 已配置, 初始化 MySQL 客户端
	var err error
	if MySQL, err = gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}); err != nil {
		return err
	}
	// 已配置, 初始化 Oracle 客户端
	if Oracle, err = gorm.Open(oracle.Open(os.Getenv("ORACLE_DSN")), &gorm.Config{}); err != nil {
		return err
	}
	// 初始化成功后，判断是否需要进行数据库表结构迁移
	if len(config.Version) != 31 {
		logger.Log("应用版本号有误, 无法计算构建时间, 开始进行数据库结构自动迁移.")
		autoMigrate()
	} else {
		// 通过版本号, 计算构建时间
		builtTime, _ := time.ParseInLocation("2006/01/02 15:04:05", config.Version[10:29], time.Local) // 使用 parseInLocation 将字符串格式化返回本地时区时间
		// 判断构建时间是否超过一小时
		if time.Since(builtTime).Hours() > 1 {
			// 超过一小时, 不需要进行迁移
			logger.Log("应用构建于 " + fmt.Sprint(time.Since(builtTime)) + " 前, 跳过数据库结构自动迁移.")
		} else {
			// 没有超过一小时, 进行表结构自动迁移
			logger.Log("应用构建时间没有超过一小时, 开始进行数据库结构自动迁移.")
			autoMigrate()
		}
	}
	return nil
}

// autoMigrate 可以完成表结构自动迁移
func autoMigrate() {
	MySQL.AutoMigrate(
		// Log 日志
		structs.Log{},
		// News 资讯
		structs.NewsTopic{},
		structs.NewsNews{},
		// 用户
		structs.User{},
		structs.UserLoginFailLog{},
		// 电话簿
		structs.TelephoneBookGroup{},
		structs.TelephoneBookContacts{},
	)
	// 判断运行环境, 如果是 release 则初始化仅在生产环境部署的表
	if os.Getenv("GIN_MODE") == "release" {
		logger.Log("应用暂不存在仅在生产环境部署的表")
	}
}
