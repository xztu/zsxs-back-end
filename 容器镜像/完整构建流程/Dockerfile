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

# 使用 alpine 作为应用的基础镜像
FROM alpine

# 更换软件源 # https://www.songma.com/news/txtlist_i29123v.html
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 修正时区
# RUN apk add --no-cache tzdata
# ENV TZ=Asia/Shanghai
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apk update && apk add libaio libnsl libc6-compat gcc

COPY instantclient-basiclite-linux.x64-12.2.0.1.0.zip ./
RUN unzip instantclient-basiclite-linux.x64-12.2.0.1.0.zip && \
    mv instantclient_12_2/ /usr/lib/ && \
    rm instantclient-basiclite-linux.x64-12.2.0.1.0.zip && \
    # ln /usr/lib/instantclient_12_2/libclntsh.so.12.1 /usr/lib/libclntsh.so && \
    # ln /usr/lib/instantclient_12_2/libmql1.so /usr/lib/libmql1.so && \
    # ln /usr/lib/instantclient_12_2/libipc1.so /usr/lib/libipc1.so && \
    # ln /usr/lib/instantclient_12_2/libnnz12.so /usr/lib/libnnz12.so && \
    # ln /usr/lib/instantclient_12_2/libons.so /usr/lib/libons.so && \
    # ln /usr/lib/instantclient_12_2/libclntshcore.so.12.1 /usr/lib/libclntshcore.so.12.1 && \
    # 解决方法来自 https://helloworldkb.com/16357096/在Alpine-Linux上的Oracle
    # libnsl
    ln /usr/lib/libnsl.so.2 /usr/lib/libnsl.so.1 && \
    # libc6-compat
    ln /lib/libc.so.6 /usr/lib/libresolv.so.2 && \
    ln /lib64/ld-linux-x86-64.so.2 /usr/lib/ld-linux-x86-64.so.2

# 配置客户端需要的环境变量
ENV LD_LIBRARY_PATH /usr/lib/instantclient_12_2
ENV ORACLE_HOME /usr/lib/instantclient_12_2

# FIX: not found
# BY: https://my.oschina.net/u/727875/blog/4556125
# RUN mkdir /lib64 && \
    # ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# 上文中安装某依赖后便不需要这个操作了

# 将构建好的二进制文件拷贝进镜像
COPY --from=builder /app/server /app/server

# 暴露端口
EXPOSE 8080

# 启动 Web 服务
CMD ["/app/server"]
