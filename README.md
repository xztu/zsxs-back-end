# 掌上忻师后端服务
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/xztu/zsxs-back-end)](https://goreportcard.com/report/github.com/xztu/zsxs-back-end) [![单元测试](https://github.com/xztu/zsxs-back-end/actions/workflows/unit-test.yml/badge.svg)](https://github.com/xztu/zsxs-back-end/actions/workflows/unit-test.yml) [![代码质量分析](https://github.com/xztu/zsxs-back-end/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/xztu/zsxs-back-end/actions/workflows/codeql-analysis.yml) [![codecov](https://codecov.io/gh/xztu/zsxs-back-end/branch/main/graph/badge.svg)](https://codecov.io/gh/xztu/zsxs-back-end)  
基于微信小程序云托管构建的掌上忻师微信小程序后端服务

## 待办列表
- [x] 首页 `/index`
    - [x] 教务 `/jw`
        - [x] 课表 `/timetable`
            - [x] 获取单日课表 `/single/:Day`
                - 按照 OpenID 获取登陆账号进行查找
            - [x] 获取整周的课表 `/all?username=xxx`
                - 配置了学号参数, 按照学号进行查找, 并返回所属人姓名
                - 没有配置学号参数, 按照 OpenID 获取登陆账号进行查找
        - [x] 成绩 `/score`
            - [x] 期末成绩 ( 获取上一个学期的期末成绩 ) `/last-semester?username=xxx`
                - 配置了学号参数, 按照学号进行查找, 并返回所属人姓名
                - 没有配置学号参数, 按照 OpenID 获取登陆账号进行查找
            - [x] 查询成绩 ( 全部成绩 ) `/all?username=x`
                - 配置了学号参数, 按照学号进行查找, 并返回所属人姓名
                - 没有配置学号参数, 按照 OpenID 获取登陆账号进行查找
            - [x] 课程成绩录入密码 `/entry-password`
                - 按照 OpenID 获取登陆账号进行查找
                - 教师类用户专用
    - [x] 电话簿 `/telephone-book`
        - [x] 分页获取分组 `/group?page=1&limit=10`
        - [x] 分页获取联系人 `/contacts?page=1&limit=10`
        - [x] 搜索 `/search?name=xx`
        - [x] 获取详情 `/detail?cid=1`
- [x] 资讯 `/news`
    - [x] 获取话题列表 `/news/topics/list`
    - [x] 获取某话题内的资讯列表 `/news/news/:Topic/list?limit=x&page=x`
- [x] 更多 `/more`
    - [x] 账号 `/account`
        - [x] 获取登陆状态 `/info`
            - 返回账号、用户类别
        - [x] 登陆 `/sign-in`
        - [x] 退出登陆 `/sign-out`

## 版本号命名规则
Git Short Commit ID [ 构建时间 ]

> 如: 0000000 [ 2021/02/19 13:58:00 ]  
> 代表: Git 提交 ID 为 0000000，构建时间是 2021年02月19日13时58分00秒。

## 部署指南
1. 云托管部署运行:
    - 方式请参考：[腾讯云 云托管 ( Tencent CloudBase Run ) 部署指南](https://cloud.tencent.com/document/product/1243/49235)  
1. 原生方式运行:
    - 编译或交叉编译源码
    - 配置环境变量
    - 启动程序
1. 容器方式运行:
    - 使用 Dockerfile 构建镜像
    - 启动镜像, 在启动参数中配置需要的环境变量
> 环境变量:
> | 变量名 | 用途 | 类型 | 示例 | 是否必选 |
> | - | - | - | - | - |
> | MYSQL_DSN | MySQL 数据源名称 | 字符串 | 用户名:密码@tcp(IP地址:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local | 必须  |
> | ORACLE_DSN | Oracle 数据源名称 | 字符串 | 用户名/密码@IP地址:端口/服务名 | 必须  |
> | GIN_MODE | Gin 运行模式, 同时被用作项目调试模式的配置项 | 字符串 | release or 其他内容 | 可选  |

## 开源声明
	Given enough eyeballs, all bugs are shallow. ( @Linus Torvalds )
	曝光足够，所有的 Bug 都是显而易见的。( "Linus法则" @Linux之父 Linus Torvalds )
我们相信开源本身所蕴含的的开放、协作与自由的精神将会给本项目注入更多新鲜的血液。  
原则上，您可以在遵循当地法律的情况下无偿地将本项目中包含的代码或完整副本应用于各种非商业用途。  
*但是我们仍旧建议您阅读并遵循以下条款：*
> [MIT LICENSE](https://github.com/xztu/zsxs-back-end/blob/main/LICENSE)
