/*
   @Time : 2021/2/22 6:52 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : timetable
   @Description: 课表
*/

package jw

import (
	"github.com/gin-gonic/gin"
	"github.com/xztu/zsxs-back-end/commons/database/orm"
	"github.com/xztu/zsxs-back-end/commons/database/structs"
	"github.com/xztu/zsxs-back-end/commons/logger"
	"github.com/xztu/zsxs-back-end/commons/response"
	"github.com/xztu/zsxs-back-end/commons/utils"
	"net/http"
	"strings"
	"time"
)

// 学期信息
type semesterInfo struct {
	Year     string `gorm:"column:XN"` // 学年
	Semester string `gorm:"column:XQ"` // 学期
}

// 课表节点
type timetable struct {
	Number   string `gorm:"column:DJJ"` // 第几节
	Time     string // 上课时间
	Name     string `gorm:"column:KCMC"` // 课程名称
	Location string `gorm:"column:SKDD"` // 上课地点
}

// 课表
type timetables struct {
	Role         string              // 角色
	Username     string              // 用户名
	RealName     string              // 真实姓名
	SemesterInfo semesterInfo        // 学期信息
	Timetables   map[int][]timetable // 课表信息
}

// TimetableGetSingle 课表 获取单日课表
func TimetableGetSingle(c *gin.Context) {
	c.Set("APIType", "TimetableGetSingle") // 添加接口类型到上下文中

	// 检查是否传入 OPENID
	if c.GetHeader("X-WX-OPENID") == "" {
		c.JSON(http.StatusBadRequest, response.Message("未传入 OPENID"))
		return
	}

	// 校验该 OPENID 是否已经登陆某账户
	userInfo := structs.User{}
	orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
	if userInfo.ID == 0 {
		c.JSON(http.StatusForbidden, response.Message("未登陆"))
		return
	}

	// 保存基本信息
	timetablesInfo := timetables{Role: userInfo.Role, Username: userInfo.Username, RealName: userInfo.RealName, Timetables: map[int][]timetable{}}

	// 根据用户类型执行查询
	if userInfo.Role == "学生" {
		// 查询学生课表
		// 查询学年学期信息
		tempSemesterInfo := semesterInfo{}
		orm.Oracle.Raw("SELECT XN, XQ FROM XSKCB WHERE XH = ? AND ROWNUM = 1 ORDER BY XKKH DESC", userInfo.Username).Scan(&tempSemesterInfo)
		timetablesInfo.SemesterInfo = tempSemesterInfo
		logger.DebugToJson("学年学期信息", tempSemesterInfo)
		// 查询课表
		var timetablesList []timetable
		orm.Oracle.Raw("SELECT DJJ, SUBSTR(KCB,1,INSTR(KCB,'<br>',1,1)-1) AS KCMC,SUBSTR(KCB,INSTR(KCB,'<br>',1,3)+4,40) || ' ( ' || SUBSTR(KCB,INSTR(KCB,'<br>',1,2)+4,INSTR(KCB,'<br>',1,3)-(INSTR(KCB,'<br>',1,2)+4)) || ' )' AS SKDD FROM XSKCB WHERE XN = ? AND XQ = ? AND XH = ? AND XQJ = ? ORDER BY DJJ ASC", timetablesInfo.SemesterInfo.Year, timetablesInfo.SemesterInfo.Semester, userInfo.Username, utils.StringToInt(c.Param("Day"))).Scan(&timetablesList)
		timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))] = timetablesList
	} else {
		// 查询教师课表
		weeks := []string{"", "%周一%", "%周二%", "%周三%", "%周四%", "%周五%", "%周六%", "%周日%"}
		// 查询学年学期信息
		tempSemesterInfo := semesterInfo{}
		orm.Oracle.Raw("SELECT XN,XQ FROM XXMC").Scan(&tempSemesterInfo)
		timetablesInfo.SemesterInfo = tempSemesterInfo
		logger.DebugToJson("学年学期信息", tempSemesterInfo)
		// 查询课表
		var timetablesList []timetable
		orm.Oracle.Raw("SELECT JXRW.KCMC, JXRW.SKSJ AS DJJ, JXRW.SKDD || ' ( ' || (SELECT BJHZ FROM JXBMCVIEW WHERE XKKH=JXRW.XKKH) || ', ' || JXRW.RS || '人 )' AS SKDD FROM (SELECT XKKH, KCMC, SKSJ, SKDD, SUM(RS) AS RS FROM JXRWB WHERE SKSJ LIKE ? AND XKKH LIKE ? AND JSZGH = ? AND SKSJ IS NOT NULL GROUP BY XKKH,KCMC,SKSJ,SKDD) JXRW", weeks[utils.StringToInt(c.Param("Day"))], "("+timetablesInfo.SemesterInfo.Year+"-"+timetablesInfo.SemesterInfo.Semester+")-%", userInfo.Username).Scan(&timetablesList)
		logger.DebugToJson("timetablesList", timetablesList)
		// 重新分割课程表
		for _, timetableInfo := range timetablesList {
			tempNumber := strings.Split(timetableInfo.Number, ";")
			tempLocation := strings.Split(timetableInfo.Location, ";")
			for index, location := range tempLocation {
				tempNumberToSplit := strings.Split(tempNumber[index], "节{")
				tempNumberToSplit = strings.Split(tempNumberToSplit[0][9:], ",")
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))] = append(timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))], timetable{
					Number:   tempNumberToSplit[0],
					Name:     timetableInfo.Name,
					Location: location,
				})
			}
		}
	}

	logger.DebugToJson("基本信息", timetablesInfo)

	// 替换上课节数信息，并增加上课时间
	for timetablesIndexForRange, timetablesInfoForRange := range timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))] {
		switch timetablesInfoForRange.Number {
		case "1":
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Number = "1 ~ 2"
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "08:00 - 09:30"
		case "3":
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Number = "3 ~ 4"
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "09:45 - 11:15"
		case "5":
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "11:25 - 12:10"
		case "6":
			if time.Now().Month() < 5 || time.Now().Month() > 9 {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "14:30 - 15:15"
			} else {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "15:00 - 15:45"
			}
		case "7":
			if time.Now().Month() < 5 || time.Now().Month() > 9 {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "15:25 - 16:10"
			} else {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "15:55 - 16:40"
			}
		case "8":
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Number = "8 ~ 9"
			if time.Now().Month() < 5 || time.Now().Month() > 9 {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "16:20 - 17:50"
			} else {
				timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "16:50 - 18:20"
			}
		case "10":
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Number = "10 ~ 11"
			timetablesInfo.Timetables[utils.StringToInt(c.Param("Day"))][timetablesIndexForRange].Time = "19:30 - 21:00"
		}
	}

	// 返回数据
	c.JSON(http.StatusOK, response.Data(timetablesInfo))
}

// TimetableGetFull 课表 获取整周的课表
func TimetableGetFull(c *gin.Context) {
	c.Set("APIType", "TimetableGetFull") // 添加接口类型到上下文中

	// 保存基本信息并初始化课表 Map
	timetablesInfo := timetables{Username: c.Query("username"), Timetables: map[int][]timetable{}}

	if timetablesInfo.Username != "" {
		// 配置了查询用户名
		// 提交请求前检查参数是否合法
		if !utils.CheckSqlInject(timetablesInfo.Username) {
			c.JSON(http.StatusForbidden, response.Message("参数不合法, 存在 SQL 注入风险"))
			return
		}
		// 使用账号取出加密后的密码
		type Result struct {
			Role     string
			RealName string `gorm:"column:XM"`
			Password string
		}
		var result Result
		if utils.IsNum(timetablesInfo.Username) && len(timetablesInfo.Username) == 12 {
			// 学生账号
			orm.Oracle.Raw("SELECT '学生' AS Role, XM, MM AS Password FROM xsjbxxb where XH = ?", timetablesInfo.Username).Scan(&result)
		} else {
			// 教师账号
			orm.Oracle.Raw("SELECT JS AS Role, XM, KL AS Password FROM yhb where YHM = ?", timetablesInfo.Username).Scan(&result)
		}
		// 用户不存在
		if result.Role == "" {
			c.JSON(http.StatusForbidden, response.Message("用户不存在"))
			return
		}
		// 添加用户信息
		timetablesInfo.Role = result.Role
		timetablesInfo.RealName = result.RealName
	} else {
		// 没有配置查询用户名
		// 检查是否传入 OPENID
		if c.GetHeader("X-WX-OPENID") == "" {
			c.JSON(http.StatusBadRequest, response.Message("未传入 OPENID"))
			return
		}

		// 校验该 OPENID 是否已经登陆某账户
		userInfo := structs.User{}
		orm.MySQL.Where("open_id = ?", c.GetHeader("X-WX-OPENID")).Find(&userInfo)
		if userInfo.ID == 0 {
			c.JSON(http.StatusForbidden, response.Message("未登陆"))
			return
		}

		// 添加用户信息
		timetablesInfo.Role = userInfo.Role
		timetablesInfo.Username = userInfo.Username
		timetablesInfo.RealName = userInfo.RealName
	}

	// 根据用户类型执行查询
	if timetablesInfo.Role == "学生" {
		// 查询学生课表
		// 查询学年学期信息
		tempSemesterInfo := semesterInfo{}
		orm.Oracle.Raw("SELECT XN, XQ FROM XSKCB WHERE XH = ? AND ROWNUM = 1 ORDER BY XKKH DESC", timetablesInfo.Username).Scan(&tempSemesterInfo)
		timetablesInfo.SemesterInfo = tempSemesterInfo
		logger.DebugToJson("学年学期信息", tempSemesterInfo)
		// 查询课表
		var timetablesList []struct {
			Day      int    `gorm:"column:XQJ"`
			Number   string `gorm:"column:DJJ"`  // 第几节
			Name     string `gorm:"column:KCMC"` // 课程名称
			Location string `gorm:"column:SKDD"` // 上课地点
		}
		orm.Oracle.Raw("SELECT XQJ, DJJ, SUBSTR(KCB,1,INSTR(KCB,'<br>',1,1)-1) AS KCMC,SUBSTR(KCB,INSTR(KCB,'<br>',1,3)+4,40) || ' ( ' || SUBSTR(KCB,INSTR(KCB,'<br>',1,2)+4,INSTR(KCB,'<br>',1,3)-(INSTR(KCB,'<br>',1,2)+4)) || ' )' AS SKDD FROM XSKCB WHERE XN = ? AND XQ = ? AND XH = ? ORDER BY DJJ ASC", timetablesInfo.SemesterInfo.Year, timetablesInfo.SemesterInfo.Semester, timetablesInfo.Username).Scan(&timetablesList)
		for _, timetablesListInfo := range timetablesList {
			timetablesInfo.Timetables[timetablesListInfo.Day] = append(timetablesInfo.Timetables[timetablesListInfo.Day], timetable{
				Number:   timetablesListInfo.Number,
				Name:     timetablesListInfo.Name,
				Location: timetablesListInfo.Location,
			})
		}
	} else {
		// 查询学年学期信息
		tempSemesterInfo := semesterInfo{}
		orm.Oracle.Raw("SELECT XN,XQ FROM XXMC").Scan(&tempSemesterInfo)
		timetablesInfo.SemesterInfo = tempSemesterInfo
		logger.DebugToJson("学年学期信息", tempSemesterInfo)
		// 查询课表
		var timetablesList []timetable
		orm.Oracle.Debug().Raw("SELECT JXRW.KCMC, JXRW.SKSJ AS DJJ, JXRW.SKDD || ' ( ' || (SELECT BJHZ FROM JXBMCVIEW WHERE XKKH=JXRW.XKKH) || ', ' || JXRW.RS || '人 )' AS SKDD FROM (SELECT XKKH, KCMC, SKSJ, SKDD, SUM(RS) AS RS FROM JXRWB WHERE XKKH LIKE ? AND JSZGH = ? AND SKSJ IS NOT NULL GROUP BY XKKH,KCMC,SKSJ,SKDD) JXRW", "("+timetablesInfo.SemesterInfo.Year+"-"+timetablesInfo.SemesterInfo.Semester+")-%", timetablesInfo.Username).Scan(&timetablesList)
		logger.DebugToJson("timetablesList", timetablesList)
		// 重新分割课程表
		for _, timetableInfo := range timetablesList {
			tempNumber := strings.Split(timetableInfo.Number, ";")
			tempLocation := strings.Split(timetableInfo.Location, ";")
			for index, location := range tempLocation {
				tempNumberToSplit := strings.Split(tempNumber[index], "节{")
				tempNumberToSplit = strings.Split(tempNumberToSplit[0][9:], ",")
				switch tempNumber[index][3:6] {
				case "一":
					timetablesInfo.Timetables[1] = append(timetablesInfo.Timetables[1], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "二":
					timetablesInfo.Timetables[2] = append(timetablesInfo.Timetables[2], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "三":
					timetablesInfo.Timetables[3] = append(timetablesInfo.Timetables[3], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "四":
					timetablesInfo.Timetables[4] = append(timetablesInfo.Timetables[4], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "五":
					timetablesInfo.Timetables[5] = append(timetablesInfo.Timetables[5], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "六":
					timetablesInfo.Timetables[6] = append(timetablesInfo.Timetables[6], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				case "日":
					timetablesInfo.Timetables[7] = append(timetablesInfo.Timetables[7], timetable{Number: tempNumberToSplit[0], Name: timetableInfo.Name, Location: location})
				}
			}
		}
	}

	logger.DebugToJson("基本信息", timetablesInfo)

	// 替换上课节数信息，并增加上课时间
	for timetablesIndexForRangeDay := range timetablesInfo.Timetables {
		for timetablesIndexForRange, timetablesInfoForRange := range timetablesInfo.Timetables[timetablesIndexForRangeDay] {
			switch timetablesInfoForRange.Number {
			case "1":
				timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "08:00 - 09:30"
			case "3":
				timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "09:45 - 11:15"
			case "5":
				timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "11:25 - 12:10"
			case "6":
				if time.Now().Month() < 5 || time.Now().Month() > 9 {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "14:30 - 15:15"
				} else {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "15:00 - 15:45"
				}
			case "7":
				if time.Now().Month() < 5 || time.Now().Month() > 9 {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "15:25 - 16:10"
				} else {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "15:55 - 16:40"
				}
			case "8":
				if time.Now().Month() < 5 || time.Now().Month() > 9 {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "16:20 - 17:50"
				} else {
					timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "16:50 - 18:20"
				}
			case "10":
				timetablesInfo.Timetables[timetablesIndexForRangeDay][timetablesIndexForRange].Time = "19:30 - 21:00"
			}
		}
	}

	// 返回数据
	c.JSON(http.StatusOK, response.Data(timetablesInfo))
}
