#build stage
FROM golang:alpine AS builder
# RUN apk add --no-cache git
WORKDIR /go/src/app
COPY golang_main .
ENV GOPROXY=https://goproxy.io
RUN go get github.com/docker/docker/api/types
RUN go get github.com/docker/docker/client
RUN go get github.com/lib/pq
RUN go get github.com/kataras/iris/v12
RUN go build -v -o main 

#final stage
FROM alpine:latest
LABEL Name=smartcodeql Version=0.0.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app .

ENTRYPOINT ./main
EXPOSE 3000
