# 发布包里的 copysync-server 是静态链接的 Linux 二进制，直接放进空镜像即可。
#
#   docker compose up -d
#
# 也可以从源码构建：在仓库根目录执行
#   docker build -f server/deploy/Dockerfile.source -t copysync-server .
FROM scratch
COPY copysync-server /copysync-server
EXPOSE 8787/tcp 3478/udp
ENTRYPOINT ["/copysync-server"]
