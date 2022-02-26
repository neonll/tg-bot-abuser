FROM alpine:3.14
RUN apk add --no-cache tzdata
ADD bin/tg_bot_abuser_arm64 app/app
ENTRYPOINT ["app/app", "/conf/config.yaml"]