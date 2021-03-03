/*
   @Time : 2021/3/3 11:04 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : utils_test
   @Description: 单元测试 工具
*/

package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestCheckSqlInject 测试 CheckSqlInject 是否可以进行 SQL 注入过滤
func TestCheckSqlInject(t *testing.T) {
	Convey("测试 CheckSqlInject 是否可以进行 SQL 注入过滤", t, func() {
		So(CheckSqlInject("普通字符串"), ShouldBeTrue)
		So(CheckSqlInject("注入字符串' and"), ShouldBeFalse)
	})
}

// TestStringToInt 测试 StringToInt 是否可以将字符串转换为整型
func TestStringToInt(t *testing.T) {
	Convey("测试 StringToInt 是否可以将字符串转换为整型", t, func() {
		So(StringToInt("a"), ShouldEqual, 0)
		So(StringToInt("a1"), ShouldEqual, 0)
		So(StringToInt("123"), ShouldEqual, 123)
	})
}

// TestReverseString 测试 ReverseString 是否可以反转字符串
func TestReverseString(t *testing.T) {
	Convey("测试 ReverseString 是否可以反转字符串", t, func() {
		So(ReverseString("1234"), ShouldEqual, "4321")
	})
}

// TestIsNum 测试 IsNum 是否可以判断字符串是否为数字
func TestIsNum(t *testing.T) {
	Convey("测试 IsNum 是否可以判断字符串是否为数字", t, func() {
		So(IsNum("a"), ShouldBeFalse)
		So(IsNum("a1"), ShouldBeFalse)
		So(IsNum("1"), ShouldBeTrue)
		So(IsNum("1.1"), ShouldBeTrue)
	})
}
