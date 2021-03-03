/*
   @Time : 2021/2/20 11:04 上午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : logger
   @Description: 日志组件
*/

package logger

import (
	"encoding/json"
	"fmt"
	"github.com/xztu/zsxs-back-end/commons/config"
	"io"
	"os"
	"runtime/debug"
	"time"
)

// DefaultWriter 用于日志输出的 io.Writer
var DefaultWriter io.Writer = os.Stdout

// Log 打印日志
func Log(log string) {
	_, _ = fmt.Fprintf(DefaultWriter, "[ 日志 ] %v | %s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		log,
	)
}

// Error 打印错误日志及调用堆栈
func Error(err error) {
	_, _ = fmt.Fprintf(DefaultWriter, "[ 错误 ] %v | %s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		err,
	)
	_, _ = fmt.Fprintf(DefaultWriter, "[ 错误 - 调用堆栈 ] %v\n%s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		debug.Stack(), // 获取调用堆栈
	)
}

// Panic 打印错误后抛出 PANIC
func Panic(err error) {
	_, _ = fmt.Fprintf(DefaultWriter, "[ 异常 - PANIC ] %v | %s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		err,
	)
	_, _ = fmt.Fprintf(DefaultWriter, "[ 异常 - PANIC - 调用堆栈 ] %v\n%s\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		debug.Stack(), // 输出调用堆栈
	)
	panic(err)
}

// DebugToJson 在调试模式开启时将参数调试输出为 Json 字符串
func DebugToJson(key string, value interface{}) {
	if config.IsDebugMode() {
		jsonStrings, _ := json.Marshal(value)
		_, _ = fmt.Fprintf(DefaultWriter, "[ 调试 - JSON ] %v | %s --> %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			key,
			jsonStrings,
		)
	}
}

// DebugToString 在调试模式开启时将参数调试输出为字符串
func DebugToString(key string, value interface{}) {
	if config.IsDebugMode() {
		_, _ = fmt.Fprintf(DefaultWriter, "[ 调试 - 字符串 ] %v | %s --> %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			key,
			fmt.Sprint(value),
		)
	}
}
