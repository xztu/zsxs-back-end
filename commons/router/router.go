/*
   @Time : 2021/2/20 11:21 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : router
   @Description: 路由及路由相关的业务
*/

package router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/xztu/zsxs-back-end/commons/config"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/handlers/index/jw"
	"github.com/xztu/zsxs-back-end/handlers/index/telephone_book"
	"github.com/xztu/zsxs-back-end/handlers/more"
	"github.com/xztu/zsxs-back-end/handlers/news"
	"io/ioutil"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 关闭 Gin 的控制台彩色输出
	gin.DisableConsoleColor()

	// 使用默认配置初始化路由
	router := gin.Default()

	// 添加版本号
	router.Use(func(c *gin.Context) {
		c.Header("Server", "ZSXS - "+config.Version)
	})

	// 添加记录日志的中间件
	router.Use(func(c *gin.Context) {
		c.Set("APIType", "Unknown")                              // 设置接口类型为未知
		body, _ := c.GetRawData()                                // 读出请求 Body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 把读过的字节流重新放到 body
		c.Next()                                                 // 请求处理结束后执行
		orm.MySQL.Create(&structs.Log{
			Source:             c.GetHeader("X-WX-SOURCE"),
			IP:                 c.ClientIP(),
			APIType:            c.GetString("APIType"),
			OpenID:             c.GetHeader("X-WX-OPENID"),
			Method:             c.Request.Method,
			URL:                c.Request.URL.String(),
			Body:               string(body),
			ResponseStatusCode: c.Writer.Status(),
		})
	})

	// 首页
	indexGroup := router.Group("/index")
	{
		// 教务
		jwGroup := indexGroup.Group("/jw")
		{
			// 课表
			timetableGroup := jwGroup.Group("/timetable")
			{
				// 获取单日课表
				timetableGroup.GET("/week/:Day", jw.TimetableGetSingle)

				// 获取整周的课表
				timetableGroup.GET("/all", jw.TimetableGetFull)
			}

			// 成绩
			scoreGroup := jwGroup.Group("/score")
			{
				// 期末成绩 ( 获取上一个学期的期末成绩 )
				scoreGroup.GET("/last-semester", jw.ScoreLastSemester)

				// 查询成绩 ( 全部成绩 )
				scoreGroup.GET("/all", jw.ScoreAll)

				// 查询课程成绩录入密码
				scoreGroup.GET("/entry-password", jw.ScoreEntryPassword)
			}
		}

		// 电话簿
		telephoneBookGroup := indexGroup.Group("/telephone-book")
		{
			// 分页获取分组
			telephoneBookGroup.GET("/group", telephoneBook.PaginationGetGroup)

			// 分页获取联系人
			telephoneBookGroup.GET("/contacts", telephoneBook.PaginationGetContacts)

			// 使用姓名模糊搜索联系人
			telephoneBookGroup.GET("/fuzzy-search/:Name", telephoneBook.FuzzySearchContactsByName)

			// 获取详情
			telephoneBookGroup.GET("/detail", telephoneBook.GetContactsInfo)
		}
	}

	// 资讯
	newsGroup := router.Group("/news")
	{
		// 获取话题列表
		newsGroup.GET("/topics/list", news.GetTopics)

		// 获取某话题内的资讯列表
		newsGroup.GET("/news/:Topic/list", news.GetNews)
	}

	// 更多
	moreGroup := router.Group("/more")
	{
		// 账号
		accountGroup := moreGroup.Group("/account")
		{
			// 登陆
			accountGroup.POST("/sign-in", more.AccountSignIn)

			// 退出登陆
			accountGroup.POST("/sign-out", more.AccountSignOut)

			// 获取登陆状态
			accountGroup.GET("/info", more.AccountGetInfo)
		}
	}

	return router
}
