FROM alpine:20251224

# 建议使用 LABEL 替代过时的 MAINTAINER 指令
LABEL maintainer="WesChen(chenxuhua0530@163.com)"

# 定义工作目录（/app 更加简洁直观）
WORKDIR /app

# 将二进制文件复制到当前目录
COPY cf-license-server-v0.0.8-linux-amd64 ./cf-license-server

# 赋予执行权限（非常重要，否则 CMD 可能会报 Permission Denied）
RUN chmod +x ./cf-license-server


# 暴露端口
EXPOSE 8080

# 使用数组形式启动，并建议写相对路径或确保在 PATH 中
# flag 包需要参数分开传递
CMD ["/app/cf-license-server", "--port", "8080"]