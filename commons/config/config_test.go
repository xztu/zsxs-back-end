/*
   @Time : 2021/3/3 9:59 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : config_test
   @Description: 单元测试 配置组件
*/

package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

// TestIsDebugMode 测试 IsDebugMode 能不能获取到是否运行调试模式
func TestIsDebugMode(t *testing.T) {
	Convey("测试 IsDebugMode 能不能获取到是否运行调试模式", t, func() {
		// 未配置 GIN_MODE 时, 返回 true
		So(IsDebugMode(), ShouldBeTrue)

		// 配置 GIN_MODE = release 时, 返回 false
		os.Setenv("GIN_MODE", "release")
		So(IsDebugMode(), ShouldBeFalse)

		// 配置 GIN_MODE != release 时, 返回 true
		os.Setenv("GIN_MODE", "test")
		So(IsDebugMode(), ShouldBeTrue)
	})
}
