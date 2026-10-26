FROM alpine:3.22

ENV TZ=Asia/Shanghai

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

# deploy.sh 将二进制、启动脚本和配置组装到 release/
ADD release /usr/local/src/generalgateway/

WORKDIR /usr/local/src/generalgateway
RUN chmod +x ./start.sh ./generalgateway

USER app

EXPOSE 8080

ENTRYPOINT ["./start.sh"]
