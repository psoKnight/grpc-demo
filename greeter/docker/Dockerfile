FROM alpine:3.22

ENV TZ=Asia/Shanghai

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

# deploy.sh 将二进制、启动脚本和配置组装到 release/
ADD release /usr/local/src/greeter/

WORKDIR /usr/local/src/greeter
RUN chmod +x ./start.sh ./greeter

USER app

EXPOSE 50051 50052

ENTRYPOINT ["./start.sh"]
