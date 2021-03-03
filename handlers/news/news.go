/*
   @Time : 2021/2/20 11:29 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : news
   @Description: 资讯模块的接口
*/

package news

import (
	"github.com/gin-gonic/gin"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/commons/response"
	"github.com/xztu/zsxs-back-end/commons/utils"
	"net/http"
	"time"
)

// GetTopics 获取话题列表
func GetTopics(c *gin.Context) {
	c.Set("APIType", "NewsGetTopics") // 添加接口类型到上下文中

	var info []struct {
		ID   uint   // ID
		Name string // 名称
		Icon string // 图标
	}
	orm.MySQL.Raw("SELECT id, `name`, icon FROM news_topics ORDER BY `order`, id DESC").Scan(&info)
	c.JSON(http.StatusOK, response.Data(info))
}

// GetNews 获取某话题内的资讯列表
func GetNews(c *gin.Context) {
	c.Set("APIType", "NewsGetNews") // 添加接口类型到上下文中
	page := 0
	if utils.StringToInt(c.Query("page")) > 0 {
		page = utils.StringToInt(c.Query("page")) - 1
	}
	limit := 5
	if c.Query("limit") != "" {
		limit = utils.StringToInt(c.Query("limit"))
	}
	if c.Param("Topic") == "headlines" {
		// 头条
		var info []struct {
			Time  time.Time `gorm:"column:updated_at"` // 时间
			Title string    // 标题
			Link  string    // 链接
		}
		orm.MySQL.Raw("SELECT updated_at, title, link FROM news_news WHERE headlines != 0 ORDER BY `order`, id DESC LIMIT ? OFFSET ?", limit, page*limit).Scan(&info)
		// 定义用于保存总量的变量
		total := int64(0)
		// 按照参数查询总量
		orm.MySQL.Model(structs.NewsNews{}).Where("headlines != 0").Count(&total)
		// 返回分页数据
		c.JSON(http.StatusOK, response.PaginationData(info, total))
	} else {
		// 非头条
		var info []struct {
			Time      time.Time `gorm:"column:updated_at"` // 时间
			Title     string    // 标题
			Link      string    // 链接
			Headlines bool      // 头条
		}
		orm.MySQL.Raw("SELECT updated_at, title, link, headlines FROM news_news WHERE topic_id = ? ORDER BY headlines DESC, `order`, id DESC LIMIT ? OFFSET ?", c.Param("Topic"), limit, page*limit).Scan(&info)
		// 定义用于保存总量的变量
		total := int64(0)
		// 按照参数查询总量
		orm.MySQL.Model(structs.NewsNews{}).Where("topic_id = ?", c.Param("Topic")).Count(&total)
		// 返回分页数据
		c.JSON(http.StatusOK, response.PaginationData(info, total))
	}
}
