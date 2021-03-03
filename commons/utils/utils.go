/*
   @Time : 2021/2/21 9:47 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : utils
   @Description: 工具
*/

package utils

import (
	"github.com/xztu/zsxs-back-end/commons/logger"
	"regexp"
	"strconv"
	"strings"
)

// CheckSqlInject SQL 注入过滤
// 参考了 # https://www.cnblogs.com/mafeng/p/6207988.html
func CheckSqlInject(parameter string) bool {
	// 过滤 ‘
	// ORACLE 注解 --  /**/
	// 关键字过滤 update, delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		logger.Error(err)
		return false
	}
	return !re.MatchString(strings.ToLower(parameter))
}

// StringToInt 将字符串转换为整型
// 摘自: http://www.57mz.com/programs/golang/52.html , 该文中还有 string 转 time 函数
func StringToInt(str string) int {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return i
}

// ReverseString 反转字符串
// 摘自: https://blog.csdn.net/qq_15437667/article/details/51714765
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// IsNum 判断字符串是否为数字
// 摘自: https://studygolang.com/topics/8696/comment/27119
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
