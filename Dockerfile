FROM alpine:latest

WORKDIR /app

# 复制二进制文件和配置文件
COPY imd-be .
COPY etc/config.yaml ./etc/

# 设置执行权限
RUN chmod +x imd-be

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./imd-be"]
