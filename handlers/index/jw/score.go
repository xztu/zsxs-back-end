/*
   @Time : 2021/2/22 6:55 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : score
   @Description: 成绩
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
)

// 成绩节点
type score struct {
	Name       string  // 课程名称
	Type       string  // 课程类型
	Credit     float64 // 学分
	GradePoint float64 // 绩点
	Score      float64 // 成绩
}

// 学年成绩节点
type semesterScore struct {
	SemesterInfo semesterInfo // 学期信息
	GPA          float64      // 平均学分绩点
	Scores       []score      // 成绩信息
}

// 成绩
type scores struct {
	Role          string          // 角色
	Username      string          // 用户名
	RealName      string          // 真实姓名
	GPA           float64         // 各学年平均学分绩点
	SemesterScore []semesterScore // 各学年成绩信息
}

// ScoreLastSemester 成绩 期末成绩 ( 获取上一个学期的期末成绩 )
func ScoreLastSemester(c *gin.Context) {
	c.Set("APIType", "ScoreLastSemester") // 添加接口类型到上下文中

	// 保存基本信息
	scoresInfo := scores{Username: c.Query("username")}

	if scoresInfo.Username != "" {
		// 配置了查询用户名
		// 使用账号取出加密后的密码
		type Result struct {
			Role     string
			RealName string `gorm:"column:XM"`
			Password string
		}
		var result Result
		if utils.IsNum(scoresInfo.Username) && len(scoresInfo.Username) == 12 {
			// 学生账号
			orm.Oracle.Raw("SELECT '学生' AS Role, XM, MM AS Password FROM xsjbxxb where XH = ?", scoresInfo.Username).Scan(&result)
		} else {
			// 教师账号
			c.JSON(http.StatusForbidden, response.Message("账号类型有误"))
			return
		}
		// 用户不存在
		if result.Role == "" {
			c.JSON(http.StatusForbidden, response.Message("用户不存在"))
			return
		}
		// 添加用户信息
		scoresInfo.Role = result.Role
		scoresInfo.RealName = result.RealName
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
		if userInfo.Role != "学生" {
			c.JSON(http.StatusForbidden, response.Message("账号类型有误"))
			return
		}

		// 添加用户信息
		scoresInfo.Role = userInfo.Role
		scoresInfo.Username = userInfo.Username
		scoresInfo.RealName = userInfo.RealName
	}

	// 查询学年学期信息
	tempSemesterInfo := semesterInfo{}
	orm.Oracle.Raw("SELECT XN, XQ FROM CJB WHERE XH = ? ORDER BY XKKH DESC", scoresInfo.Username).Scan(&tempSemesterInfo)
	logger.DebugToJson("学年学期信息", tempSemesterInfo)

	var tempScoresList []struct {
		Name        string  `gorm:"column:KCMC"` // 课程名称
		Type        string  `gorm:"column:KCXZ"` // 课程类型
		Credit      float64 `gorm:"column:XF"`   // 学分
		Score       float64 `gorm:"column:CJ"`   // 成绩
		RetakeScore float64 `gorm:"column:BKCJ"` // 补考成绩
		Revamp      bool    `gorm:"column:CXBJ"` // 重修标记
	}
	orm.Oracle.Raw("SELECT KCMC, XF, CJ, NVL(BKCJ,-1) AS BKCJ, CXBJ, KCXZ FROM CJB WHERE XH = ? AND XN = ? AND XQ = ?", scoresInfo.Username, tempSemesterInfo.Year, tempSemesterInfo.Semester).Scan(&tempScoresList)

	var scoresList []score

	semesterCredit := 0.0     // 学期总学分
	semesterGradePoint := 0.0 // 学期总绩点

	// 遍历成绩
	for _, scoreInfo := range tempScoresList {
		// 添加重修标记
		if scoreInfo.Revamp {
			scoreInfo.Type += " - 重修"
		}

		// 判断是否存在补考成绩
		if scoreInfo.RetakeScore != -1 {
			// 存在补考成绩
			if scoreInfo.RetakeScore > scoreInfo.Score {
				// 补考成绩大于正考成绩则使用补考成绩
				if scoreInfo.RetakeScore >= 60 {
					// 补考成绩合格则算作60分
					scoreInfo.Score = 60
				} else {
					scoreInfo.Score = scoreInfo.RetakeScore
				}
			}
		}

		// 计算绩点
		// 依据忻州师范学院本科学分制教学管理办法（http://dept.xztc.edu.cn/xsc/new/show.asp?id=165）计算
		gradePoint := 0.0
		if scoreInfo.Score < 60 {
			gradePoint = 0.0
		} else if scoreInfo.Score < 65 {
			gradePoint = 1.0
		} else if scoreInfo.Score < 70 {
			gradePoint = 1.5
		} else if scoreInfo.Score < 75 {
			gradePoint = 2.0
		} else if scoreInfo.Score < 80 {
			gradePoint = 2.5
		} else if scoreInfo.Score < 85 {
			gradePoint = 3.0
		} else if scoreInfo.Score < 90 {
			gradePoint = 3.5
		} else {
			gradePoint = 4.0
		}

		// 添加本科成绩
		scoresList = append(scoresList, score{
			Name:       scoreInfo.Name,
			Type:       scoreInfo.Type,
			Credit:     scoreInfo.Credit,
			GradePoint: gradePoint,
			Score:      scoreInfo.Score,
		})

		// 统计总学分
		semesterCredit += scoreInfo.Credit
		// 统计总学分绩点
		semesterGradePoint += scoreInfo.Credit * gradePoint
	}

	// 添加学年成绩信息
	scoresInfo.SemesterScore = append(scoresInfo.SemesterScore, semesterScore{SemesterInfo: tempSemesterInfo, GPA: semesterGradePoint / semesterCredit, Scores: scoresList})

	// 计算各学年平均学分绩点
	scoresInfo.GPA = semesterGradePoint / semesterCredit

	// 返回数据
	c.JSON(http.StatusOK, response.Data(scoresInfo))
}

// ScoreAll 成绩 查询成绩 ( 全部成绩 )
func ScoreAll(c *gin.Context) {
	c.Set("APIType", "ScoreAll") // 添加接口类型到上下文中

	// 保存基本信息
	scoresInfo := scores{Username: c.Query("username")}

	if scoresInfo.Username != "" {
		// 配置了查询用户名
		// 使用账号取出加密后的密码
		type Result struct {
			Role     string
			RealName string `gorm:"column:XM"`
			Password string
		}
		var result Result
		if utils.IsNum(scoresInfo.Username) && len(scoresInfo.Username) == 12 {
			// 学生账号
			orm.Oracle.Raw("SELECT '学生' AS Role, XM, MM AS Password FROM xsjbxxb where XH = ?", scoresInfo.Username).Scan(&result)
		} else {
			// 教师账号
			c.JSON(http.StatusForbidden, response.Message("账号类型有误"))
			return
		}
		// 用户不存在
		if result.Role == "" {
			c.JSON(http.StatusForbidden, response.Message("用户不存在"))
			return
		}
		// 添加用户信息
		scoresInfo.Role = result.Role
		scoresInfo.RealName = result.RealName
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
		if userInfo.Role != "学生" {
			c.JSON(http.StatusForbidden, response.Message("账号类型有误"))
			return
		}

		// 添加用户信息
		scoresInfo.Role = userInfo.Role
		scoresInfo.Username = userInfo.Username
		scoresInfo.RealName = userInfo.RealName
	}

	// 查询学年学期信息
	var tempSemesterInfos []semesterInfo
	orm.Oracle.Raw("SELECT XN, XQ FROM CJB WHERE XH = ? GROUP BY XN, XQ ORDER BY XN DESC, XQ DESC", scoresInfo.Username).Scan(&tempSemesterInfos)
	logger.DebugToJson("学年学期信息", tempSemesterInfos)

	allCredit := 0.0     // 学期总学分
	allGradePoint := 0.0 // 学期总绩点

	for _, tempSemesterInfo := range tempSemesterInfos {
		var tempScoresList []struct {
			Name        string  `gorm:"column:KCMC"` // 课程名称
			Type        string  `gorm:"column:KCXZ"` // 课程类型
			Credit      float64 `gorm:"column:XF"`   // 学分
			Score       float64 `gorm:"column:CJ"`   // 成绩
			RetakeScore float64 `gorm:"column:BKCJ"` // 补考成绩
			Revamp      bool    `gorm:"column:CXBJ"` // 重修标记
		}
		orm.Oracle.Raw("SELECT KCMC, XF, CJ, NVL(BKCJ,-1) AS BKCJ, CXBJ, KCXZ FROM CJB WHERE XH = ? AND XN = ? AND XQ = ?", scoresInfo.Username, tempSemesterInfo.Year, tempSemesterInfo.Semester).Scan(&tempScoresList)

		var scoresList []score

		semesterCredit := 0.0     // 学期总学分
		semesterGradePoint := 0.0 // 学期总绩点

		// 遍历成绩
		for _, scoreInfo := range tempScoresList {
			// 添加重修标记
			if scoreInfo.Revamp {
				scoreInfo.Type += " - 重修"
			}

			// 判断是否存在补考成绩
			if scoreInfo.RetakeScore != -1 {
				// 存在补考成绩
				if scoreInfo.RetakeScore > scoreInfo.Score {
					// 补考成绩大于正考成绩则使用补考成绩
					if scoreInfo.RetakeScore >= 60 {
						// 补考成绩合格则算作60分
						scoreInfo.Score = 60
					} else {
						scoreInfo.Score = scoreInfo.RetakeScore
					}
				}
			}

			// 计算绩点
			// 依据忻州师范学院本科学分制教学管理办法（http://dept.xztc.edu.cn/xsc/new/show.asp?id=165）计算
			gradePoint := 0.0
			if scoreInfo.Score < 60 {
				gradePoint = 0.0
			} else if scoreInfo.Score < 65 {
				gradePoint = 1.0
			} else if scoreInfo.Score < 70 {
				gradePoint = 1.5
			} else if scoreInfo.Score < 75 {
				gradePoint = 2.0
			} else if scoreInfo.Score < 80 {
				gradePoint = 2.5
			} else if scoreInfo.Score < 85 {
				gradePoint = 3.0
			} else if scoreInfo.Score < 90 {
				gradePoint = 3.5
			} else {
				gradePoint = 4.0
			}

			// 添加本科成绩
			scoresList = append(scoresList, score{
				Name:       scoreInfo.Name,
				Type:       scoreInfo.Type,
				Credit:     scoreInfo.Credit,
				GradePoint: gradePoint,
				Score:      scoreInfo.Score,
			})

			// 统计总学分
			semesterCredit += scoreInfo.Credit
			allCredit += scoreInfo.Credit
			// 统计总学分绩点
			semesterGradePoint += scoreInfo.Credit * gradePoint
			allGradePoint += scoreInfo.Credit * gradePoint
		}

		// 添加学年成绩信息
		scoresInfo.SemesterScore = append(scoresInfo.SemesterScore, semesterScore{SemesterInfo: tempSemesterInfo, GPA: semesterGradePoint / semesterCredit, Scores: scoresList})
	}

	// 计算各学年平均学分绩点
	scoresInfo.GPA = allGradePoint / allCredit

	// 返回数据
	c.JSON(http.StatusOK, response.Data(scoresInfo))
}

// ScoreEntryPassword 成绩 课程成绩录入密码
func ScoreEntryPassword(c *gin.Context) {
	c.Set("APIType", "ScoreEntryPassword") // 添加接口类型到上下文中

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
	if userInfo.Role == "学生" {
		c.JSON(http.StatusForbidden, response.Message("账号类型有误"))
		return
	}

	// 录入密码
	type password struct {
		Status    string `gorm:"column:LRSZ"` // 录入状态
		Name      string `gorm:"column:KCMC"` // 课程名称
		Time      string `gorm:"column:SKSJ"` // 上课时间
		Location  string `gorm:"column:SKDD"` // 上课地点
		ClassName string `gorm:"column:BJHZ"` // 班级名称
		Password  string `gorm:"column:MM"`   // 密码
	}

	// 补考成绩录入密码
	type supplementPassword struct {
		Status         string `gorm:"column:LRSZ"` // 录入状态
		Name           string `gorm:"column:KCMC"` // 课程名称
		NumberOfPeople string `gorm:"column:RS"`   // 补考人数
		Password       string `gorm:"column:MM"`   // 密码
	}

	// 信息
	type info struct {
		Role                   string               // 角色
		Username               string               // 用户名
		RealName               string               // 真实姓名
		SemesterInfo           semesterInfo         // 录入学年学期信息
		Passwords              []password           // 录入密码
		SupplementSemesterInfo semesterInfo         // 补考成绩录入学年学期信息
		SupplementPasswords    []supplementPassword // 补考成绩录入密码
	}

	// 查询成绩、补考成绩录入学年学期信息
	type tempSemesterInfoStruct struct {
		Year                      string `gorm:"column:CJXN"`     // 录入学年
		Semester                  string `gorm:"column:CJXQ"`     // 录入学期
		SupplementYearAndSemester string `gorm:"column:LRBKXNXQ"` // 补考成绩录入学年学期
	}
	tempSemesterInfo := tempSemesterInfoStruct{}
	orm.Oracle.Raw("SELECT CJXN, CJXQ, LRBKXNXQ FROM XXMC").Scan(&tempSemesterInfo)
	logger.DebugToJson("成绩、补考成绩录入学年学期信息", tempSemesterInfo)

	// 添加基本信息
	infos := info{Role: userInfo.Role, Username: userInfo.Username, RealName: userInfo.RealName, SemesterInfo: semesterInfo{Year: tempSemesterInfo.Year, Semester: tempSemesterInfo.Semester}, SupplementSemesterInfo: semesterInfo{Year: tempSemesterInfo.SupplementYearAndSemester[:9], Semester: tempSemesterInfo.SupplementYearAndSemester[9:]}}

	// 查询录入信息
	var passwords []password
	orm.Oracle.Raw("SELECT DISTINCT JXRW.LRSZ, JXRW.KCMC, JXRW.SKSJ, JXRW.SKDD, (SELECT MAX(TRIM(JXBMCVIEW.BJHZ)) FROM JXBMCVIEW WHERE JXRW.xkkh=JXBMCVIEW.XKKH) BJHZ, JXRW.MM FROM JXRWBVIEW JXRW WHERE JXRW.XKKH LIKE ? AND (JXRW.XKZT <> '4' OR JXRW.XKZT IS NULL) AND JXRW.JSZGH = ?", "("+infos.SemesterInfo.Year+"-"+infos.SemesterInfo.Semester+")-%", userInfo.Username).Scan(&passwords)
	logger.DebugToJson("录入信息", passwords)
	// 遍历信息, 进行解密
	for index := range passwords {
		passwords[index].Password = utils.DecryptPassword(passwords[index].Password)
	}
	logger.DebugToJson("解密密码后的录入信息", passwords)
	// 将信息添加到基本信息中
	infos.Passwords = passwords

	// 查询补考录入信息
	var supplementPasswords []supplementPassword
	orm.Oracle.Raw("SELECT BKMD.LRSZ, BKMD.KCMC, COUNT(*) AS RS, BKMD.MM FROM (SELECT LRSZ, KCMC, MM, XH, JSZGH FROM BKMDB WHERE BKXN = ? AND BKXQ = ? AND (SFBYBK IS NULL OR SFBYBK = '否') AND DJCBK = '1' AND BKQR = '1') BKMD, XSJBXXB WHERE BKMD.XH = XSJBXXB.xh AND BKMD.JSZGH = ? GROUP BY BKMD.LRSZ, BKMD.KCMC, BKMD.MM", infos.SupplementSemesterInfo.Year, infos.SemesterInfo.Semester, userInfo.Username).Scan(&supplementPasswords)
	logger.DebugToJson("补考录入信息", supplementPasswords)
	// 遍历信息, 进行解密
	for index := range supplementPasswords {
		supplementPasswords[index].Password = utils.DecryptPassword(supplementPasswords[index].Password)
	}
	logger.DebugToJson("解密密码后的补考录入信息", supplementPasswords)
	// 将信息添加到基本信息中
	infos.SupplementPasswords = supplementPasswords

	// 返回数据
	c.JSON(http.StatusOK, response.Data(infos))
}
