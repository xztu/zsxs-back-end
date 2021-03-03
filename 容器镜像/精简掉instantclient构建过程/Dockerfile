# 使用官方 Golang 镜像作为构建环境
FROM golang:1.15-buster as builder
WORKDIR /app
# 安装依赖
ENV GOPROXY https://goproxy.cn,direct
COPY go.* ./
RUN go mod download
# 将代码文件写入镜像
COPY . ./
# 构建二进制文件
RUN go build -v -o server

# 使用 alpine-instantclient_12_2 作为应用的基础镜像
FROM ccr.ccs.tencentyun.com/tcb-100017786839-vgjo/alpine-instantclient_12_2

# 将构建好的二进制文件拷贝进镜像
COPY --from=builder /app/server /app/server

# 暴露端口
EXPOSE 8080

# 启动 Web 服务
CMD ["/app/server"]
