/*
   @Time : 2021/2/20 11:21 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : router_test
   @Description: 单元测试 路由及路由相关的业务
*/

package router

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestInitRouter 测试 InitRouter 是否可以初始化路由
func TestInitRouter(t *testing.T) {
	Convey("测试 InitRouter 是否可以初始化路由", t, func() {
		ginEngine := InitRouter()
		So(ginEngine.AppEngine, ShouldBeFalse)

		testServer := httptest.NewServer(ginEngine)

		// 测试实际响应
		res, err := http.Get(fmt.Sprintf("%s/fake-path", testServer.URL))
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusInternalServerError)
		resp, err := ioutil.ReadAll(res.Body)
		So(err, ShouldBeNil)
		So(string(resp), ShouldBeEmpty)
	})
}
