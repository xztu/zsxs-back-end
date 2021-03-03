/*
   @Time : 2021/2/20 12:30 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : response_test
   @Description: 单元测试 响应内容
*/

package response

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// TestMessage 测试 Message 是否可以按照预期创建响应
func TestMessage(t *testing.T) {
	Convey("测试 Message 是否可以按照预期创建响应", t, func() {
		So(fmt.Sprint(Message("预期内容")), ShouldEqual, "map[Message:预期内容]")
	})
}

// TestData 测试 Data 是否可以按照预期创建响应
func TestData(t *testing.T) {
	Convey("测试 Data 是否可以按照预期创建响应", t, func() {
		So(fmt.Sprint(Data(map[string]string{"预期内容下标": "预期内容"})), ShouldEqual, "map[Data:map[预期内容下标:预期内容] Message:Success]")
	})
}

// TestPaginationData 测试 PaginationData 是否可以按照预期创建响应
func TestPaginationData(t *testing.T) {
	Convey("测试 PaginationData 是否可以按照预期创建响应", t, func() {
		So(fmt.Sprint(PaginationData(map[string]string{"预期内容下标": "预期内容"}, 6688)), ShouldEqual, "map[Data:map[预期内容下标:预期内容] Message:Success Total:6688]")
	})
}

// TestError 测试 Error 是否可以按照预期创建响应
func TestError(t *testing.T) {
	Convey("测试 Error 是否可以按照预期创建响应", t, func() {
		So(fmt.Sprint(Error("预期内容", errors.New("预期错误文本"))), ShouldEqual, "map[Error:预期错误文本 Message:预期内容]")
	})
}

// TestJsonInvalid 测试 (json) Invalid 是否可以按照预期创建响应
func TestJsonInvalid(t *testing.T) {
	Convey("测试 Error 是否可以按照预期创建响应", t, func() {
		So(fmt.Sprint(Json.Invalid(errors.New("预期错误文本"))), ShouldEqual, "map[Error:预期错误文本 Message:提交的 Json 数据不正确]")
	})
}
