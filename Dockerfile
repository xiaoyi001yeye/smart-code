#build stage
FROM golang:1.21.8-alpine3.19 AS builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY golang_main .
# ENV GOPROXY=https://goproxy.io,direct,goproxy.cn
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct,https://goproxy.io,https://mirrors.aliyun.com/goproxy/
RUN go env -w GOPRIVATE=github.com/smartcodeql
# ENV GOPRIVATE=github.com/smartcodeql
RUN go mod tidy
# RUN go get github.com/kataras/iris/v12
# RUN go get golang.org/x/xerrors
# RUN go get github.com/docker/docker/api/types
# RUN go get github.com/docker/docker/client
# RUN go get github.com/lib/pq
RUN go build -o main

#final stage
FROM alpine:latest
LABEL Name=smartcodeql Version=0.0.1
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app .

ENTRYPOINT ./main
EXPOSE 3000
