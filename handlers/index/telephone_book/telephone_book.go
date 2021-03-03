/*
   @Time : 2021/2/24 10:38 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : telephone_book
   @Description: 电话簿
*/

package telephoneBook

import (
	"github.com/gin-gonic/gin"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/commons/response"
	"github.com/xztu/zsxs-back-end/commons/utils"
	"net/http"
)

// PaginationGetGroup 分页获取分组
func PaginationGetGroup(c *gin.Context) {
	c.Set("APIType", "TelephoneBookPaginationGetGroup") // 添加接口类型到上下文中

	page := 0
	if utils.StringToInt(c.Query("page")) > 0 {
		page = utils.StringToInt(c.Query("page")) - 1
	}
	limit := 5
	if c.Query("limit") != "" {
		limit = utils.StringToInt(c.Query("limit"))
	}

	var info []struct {
		ID   uint
		Type string
		Name string
		Rank uint
	}
	orm.MySQL.Raw("SELECT id, type, `name`, rank FROM telephone_book_groups ORDER BY rank DESC LIMIT ?, ?", page*limit, limit).Scan(&info)

	total := int64(0)
	orm.MySQL.Model(structs.TelephoneBookGroup{}).Count(&total)

	c.JSON(http.StatusOK, response.PaginationData(info, total))
}

// PaginationGetContacts 分页获取联系人
func PaginationGetContacts(c *gin.Context) {
	c.Set("APIType", "TelephoneBookPaginationGetContacts") // 添加接口类型到上下文中

	if utils.StringToInt(c.Query("gid")) == 0 {
		c.JSON(http.StatusForbidden, response.Message("分组代码不合规"))
		return
	}

	page := 0
	if utils.StringToInt(c.Query("page")) > 0 {
		page = utils.StringToInt(c.Query("page")) - 1
	}
	limit := 5
	if c.Query("limit") != "" {
		limit = utils.StringToInt(c.Query("limit"))
	}

	publicOnly := false
	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		publicOnly = true
	}

	// 校验该 OPENID 是否已经登陆某账户
	userInfo := structs.User{}
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
	if userInfo.Role == "学生" {
		publicOnly = true
	}

	var info []struct {
		ID   uint
		Name string
	}
	if publicOnly {
		orm.MySQL.Raw("SELECT id, `name` FROM telephone_book_contacts WHERE group_id = ? AND public = ? LIMIT ?, ?", utils.StringToInt(c.Query("gid")), true, page*limit, limit).Scan(&info)
	} else {
		orm.MySQL.Raw("SELECT id, `name` FROM telephone_book_contacts WHERE group_id = ? LIMIT ?, ?", utils.StringToInt(c.Query("gid")), page*limit, limit).Scan(&info)
	}

	total := int64(0)
	if publicOnly {
		orm.MySQL.Model(structs.TelephoneBookContacts{}).Where("group_id = ? AND public = ?", utils.StringToInt(c.Query("gid")), true).Count(&total)
	} else {
		orm.MySQL.Model(structs.TelephoneBookContacts{}).Where("group_id = ?", utils.StringToInt(c.Query("gid"))).Count(&total)
	}

	c.JSON(http.StatusOK, response.PaginationData(info, total))
}

// FuzzySearchContactsByName 使用姓名模糊搜索联系人
func FuzzySearchContactsByName(c *gin.Context) {
	c.Set("APIType", "TelephoneBookFuzzySearchContactsByName") // 添加接口类型到上下文中

	if !utils.CheckSqlInject(c.Param("Name")) {
		c.JSON(http.StatusForbidden, response.Message("参数不合法, 存在 SQL 注入风险"))
		return
	}

	page := 0
	if utils.StringToInt(c.Query("page")) > 0 {
		page = utils.StringToInt(c.Query("page")) - 1
	}
	limit := 5
	if c.Query("limit") != "" {
		limit = utils.StringToInt(c.Query("limit"))
	}

	publicOnly := false
	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		publicOnly = true
	}

	// 校验该 OPENID 是否已经登陆某账户
	userInfo := structs.User{}
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
	if userInfo.Role == "学生" {
		publicOnly = true
	}

	var info []struct {
		ID   uint
		Name string
	}
	if publicOnly {
		orm.MySQL.Raw("SELECT id, `name` FROM telephone_book_contacts WHERE name LIKE ? AND public = ? LIMIT ?, ?", "%"+c.Param("Name")+"%", true, page*limit, limit).Scan(&info)
	} else {
		orm.MySQL.Raw("SELECT id, `name` FROM telephone_book_contacts WHERE name LIKE ? LIMIT ?, ?", "%"+c.Param("Name")+"%", page*limit, limit).Scan(&info)
	}

	total := int64(0)
	if publicOnly {
		orm.MySQL.Model(structs.TelephoneBookContacts{}).Where("name LIKE ? AND public = ?", "%"+c.Param("Name")+"%", true).Count(&total)
	} else {
		orm.MySQL.Model(structs.TelephoneBookContacts{}).Where("name LIKE ?", "%"+c.Param("Name")+"%").Count(&total)
	}

	c.JSON(http.StatusOK, response.PaginationData(info, total))
}

// GetContactsInfo 获取联系人详情
func GetContactsInfo(c *gin.Context) {
	c.Set("APIType", "TelephoneBookGetContactsInfo") // 添加接口类型到上下文中

	if utils.StringToInt(c.Query("cid")) == 0 {
		c.JSON(http.StatusForbidden, response.Message("联系人代码不合规"))
		return
	}

	publicOnly := false
	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		publicOnly = true
	}

	// 校验该 OPENID 是否已经登陆某账户
	userInfo := structs.User{}
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
	if userInfo.Role == "学生" {
		publicOnly = true
	}

	info := struct {
		GroupName         string // 所属分组
		Name              string // 姓名
		Office            string // 办公室电话
		OfficeShort       string // 办公室电话 小号
		ChinaMobile       string // 中国移动
		ChinaMobileShort  string // 中国移动 小号
		ChinaTelecom      string // 中国电信
		ChinaTelecomShort string // 中国电信 小号
		ChinaUnicom       string // 中国联通
		Remark            string // 备注
	}{}
	if publicOnly {
		orm.MySQL.Raw("SELECT telephone_book_groups.`name` AS group_name,telephone_book_contacts.`name`,telephone_book_contacts.office,telephone_book_contacts.office_short,telephone_book_contacts.china_mobile,telephone_book_contacts.china_mobile_short,telephone_book_contacts.china_telecom,telephone_book_contacts.china_telecom_short,telephone_book_contacts.china_unicom,telephone_book_contacts.remark FROM telephone_book_contacts,telephone_book_groups WHERE telephone_book_contacts.group_id=telephone_book_groups.id AND telephone_book_contacts.id=? AND telephone_book_contacts.public=?", utils.StringToInt(c.Query("cid")), true).Scan(&info)
	} else {
		orm.MySQL.Raw("SELECT telephone_book_groups.`name` AS group_name,telephone_book_contacts.`name`,telephone_book_contacts.office,telephone_book_contacts.office_short,telephone_book_contacts.china_mobile,telephone_book_contacts.china_mobile_short,telephone_book_contacts.china_telecom,telephone_book_contacts.china_telecom_short,telephone_book_contacts.china_unicom,telephone_book_contacts.remark FROM telephone_book_contacts,telephone_book_groups WHERE telephone_book_contacts.group_id=telephone_book_groups.id AND telephone_book_contacts.id=?", utils.StringToInt(c.Query("cid"))).Scan(&info)
	}

	if info.Name == "" {
		c.JSON(http.StatusOK, response.Message("未找到联系人"))
	} else {
		c.JSON(http.StatusOK, response.Data(info))
	}
}
