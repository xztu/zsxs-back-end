# 使用官方 Golang 镜像作为构建环境
FROM golang:1.15-buster

# 修正时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
