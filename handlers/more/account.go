/*
   @Time : 2021/2/21 6:04 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : account
   @Description: 账户 ( 登陆、退出登陆、获取账户信息等操作 )
*/

package more

import (
	"github.com/gin-gonic/gin"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/commons/logger"
	"github.com/xztu/zsxs-back-end/commons/response"
	"github.com/xztu/zsxs-back-end/commons/utils"
	"net/http"
	"time"
)

// 登陆
func AccountSignIn(c *gin.Context) {
	c.Set("APIType", "AccountSignIn") // 添加接口类型到上下文中

	// 绑定数据
	requestJsonMap := struct {
		Username string `json:"Username" binding:"required"` // 用户名
		Password string `json:"Password" binding:"required"` // 密码
	}{}
	// 绑定参数
	if err := c.ShouldBindJSON(&requestJsonMap); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, response.Json.Invalid(err))
		return
	}

	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		c.JSON(http.StatusBadRequest, response.Message("未传入 OPENID"))
		return
	}

	// 校验该 OPENID 是否已经登陆某账户
	userInfo := structs.User{}
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
	if userInfo.ID != 0 {
		c.JSON(http.StatusForbidden, response.Message("您已经登陆了账号"))
		return
	}

	// 获取 24h 登陆失败次数，登陆失败次数超过 5 次则拒绝登陆
	// 获取 UserLoginFailLog 中 OPENID = 要登陆的 OPENID AND CreateAt > 当前时间-1d的数据的数量
	loginFailCount := int64(0)
	orm.MySQL.Model(structs.UserLoginFailLog{}).Where("open_id = ? AND created_at > ?", c.GetHeader("X-WX-OPENID"), time.Now().AddDate(0, 0, -1)).Count(&loginFailCount)
	if loginFailCount > 5 {
		c.JSON(http.StatusForbidden, response.Message("您已经使用该微信账号在 24 小时内连续登陆失败 5 次, 已经被暂时冻结, 24 小时后自动解冻"))
		return
	}
	// 获取 UserLoginFailLog 中 Username = 要登陆的 Username AND CreateAt > 当前时间-1d的数据的数量
	orm.MySQL.Model(structs.UserLoginFailLog{}).Where("username = ? AND created_at > ?", requestJsonMap.Username, time.Now().AddDate(0, 0, -1)).Count(&loginFailCount)
	if loginFailCount > 5 {
		c.JSON(http.StatusForbidden, response.Message("您已经使用该教务系统账号在 24 小时内连续登陆失败 5 次, 已经被暂时冻结, 24 小时后自动解冻"))
		return
	}

	// 使用账号取出加密后的密码
	type Result struct {
		Role     string
		RealName string `gorm:"column:XM"`
		Password string
	}
	var result Result
	if utils.IsNum(requestJsonMap.Username) && len(requestJsonMap.Username) == 12 {
		// 学生账号
		orm.Oracle.Raw("SELECT '学生' AS Role, XM, MM AS Password FROM xsjbxxb where XH = ?", requestJsonMap.Username).Scan(&result)
	} else {
		// 提交请求前检查参数是否合法
		if !utils.CheckSqlInject(requestJsonMap.Username) {
			c.JSON(http.StatusForbidden, response.Message("用户名不合法, 存在 SQL 注入风险"))
			return
		}
		// 教师账号
		orm.Oracle.Raw("SELECT JS AS Role, XM, KL AS Password FROM yhb where YHM = ?", requestJsonMap.Username).Scan(&result)
	}

	// 用户不存在
	if result.Role == "" {
		// 保存登陆失败记录
		orm.MySQL.Create(&structs.UserLoginFailLog{
			Type:     "UserNotExist",
			OpenID:   c.GetHeader("X-WX-OPENID"),
			Username: requestJsonMap.Username,
			Password: requestJsonMap.Password,
			SourceIp: c.ClientIP(),
		})
		c.JSON(http.StatusForbidden, response.Message("用户不存在"))
		return
	}

	// 加密请求中的密码, 检查是否与记录中的密码相符
	if result.Password != utils.EncryptPassword(requestJsonMap.Password) {
		// 保存登陆失败记录
		orm.MySQL.Create(&structs.UserLoginFailLog{
			Type:     "PasswordIsIncorrect",
			OpenID:   c.GetHeader("X-WX-OPENID"),
			Username: requestJsonMap.Username,
			Password: requestJsonMap.Password,
			SourceIp: c.ClientIP(),
		})
		c.JSON(http.StatusForbidden, response.Message("密码不正确"))
		return
	}

	// 退出登陆其他用户记录
	orm.MySQL.Where("username = ?", requestJsonMap.Username).Delete(&structs.User{})

	// 创建用户记录
	orm.MySQL.Create(&structs.User{
		Role:     result.Role,
		OpenID:   c.GetHeader("X-WX-OPENID"),
		Username: requestJsonMap.Username,
		RealName: result.RealName,
	})

	// 返回成功
	c.JSON(http.StatusOK, response.Success)
}

// 退出登陆
func AccountSignOut(c *gin.Context) {
	c.Set("APIType", "AccountSignOut") // 添加接口类型到上下文中

	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		c.JSON(http.StatusBadRequest, response.Message("未传入 OPENID"))
		return
	}

	// 删除用户记录
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Delete(&structs.User{})

	// 返回成功
	c.JSON(http.StatusOK, response.Success)
}

// 获取登陆状态
func AccountGetInfo(c *gin.Context) {
	c.Set("APIType", "AccountGetInfo") // 添加接口类型到上下文中

	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		c.JSON(http.StatusBadRequest, response.Message("未传入 OPENID"))
		return
	}

	// 取出用户信息
	type Result struct {
		ID       uint
		Role     string
		Username string
		RealName string
	}
	var result Result
	orm.MySQL.Raw("SELECT id, role, username, real_name FROM users WHERE open_id = ? AND deleted_at IS NULL", c.GetHeader("X-WX-OPENID")).Scan(&result)

	// 判断用户是否登陆
	if result.ID == 0 {
		c.JSON(http.StatusBadRequest, response.Message("未登陆"))
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, response.Data(result))
}
