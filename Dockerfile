# 使用 Golang 官方镜像作为基础镜像
FROM golang:1.23.3

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件到工作目录
COPY go.mod go.sum ./

# 下载 Go 依赖
RUN go mod tidy

# 将项目源代码复制到工作目录
COPY . .

# 编译 Go 应用
RUN go build -o main .

# 容器启动时运行的命令
CMD ["./main"]
