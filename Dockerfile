FROM golang:1.14-alpine3.12 AS builder

WORKDIR /root/doghandler
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux go build -a -trimpath -o bin/doghandler cmd/doghandler/main.go


#ARG NPM_REGISTRY=https://registry.npm.taobao.org
#ARG GOPROXY="https://mirrors.aliyun.com/goproxy/"

FROM alpine:3.12

LABEL maintainer="LyleLaii"

# ENV TZ Asia/Shanghai
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories &&\
#     apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
#     && echo ${TZ} > /etc/timezone \
#     && apk del tzdata

WORKDIR /opt/doghandler

COPY --from=builder /root/doghandler/bin /opt/doghandler/bin
COPY conf/settings.yaml /etc/settings.yaml

EXPOSE 8080

ENTRYPOINT [ "bin/doghandler" ]
CMD ["-c", "/etc/settings.yaml" ]