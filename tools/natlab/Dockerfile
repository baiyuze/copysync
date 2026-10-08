# NAT 打洞实验室的运行环境：只需要 bash、iproute2 与 iptables。
FROM alpine:3.20
ARG APK_MIRROR=""
RUN if [ -n "$APK_MIRROR" ]; then \
        sed -i "s#dl-cdn.alpinelinux.org#$APK_MIRROR#g" /etc/apk/repositories; \
    fi \
    && apk add --no-cache bash iproute2 iptables
